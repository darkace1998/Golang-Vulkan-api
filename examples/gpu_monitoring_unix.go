//go:build !windows

package main

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/NVIDIA/go-nvml/pkg/nvml"
)

func (app *BenchmarkApp) initGPUMonitoring() {
	ret := nvml.Init()
	if ret != nvml.SUCCESS {
		log.Printf("Failed to initialize NVML: %v", nvml.ErrorString(ret))
		return
	}
	app.nvmlInitialized = true
	log.Println("GPU monitoring initialized")
}

func (app *BenchmarkApp) cleanupGPUMonitoring() {
	if app.nvmlInitialized {
		nvml.Shutdown()
	}
}

func (app *BenchmarkApp) getGPUStats() *GPUStats {
	// Try NVIDIA monitoring first
	if nvmlStats := app.getNvidiaGPUStats(); nvmlStats != nil {
		return nvmlStats
	}

	// Try generic Linux GPU monitoring
	if genericStats := app.getGenericGPUStats(); genericStats != nil {
		return genericStats
	}

	return nil
}

func (app *BenchmarkApp) getNvidiaGPUStats() *GPUStats {
	if !app.nvmlInitialized {
		return nil
	}

	deviceCount, ret := nvml.DeviceGetCount()
	if ret != nvml.SUCCESS || deviceCount == 0 {
		return nil
	}

	device, ret := nvml.DeviceGetHandleByIndex(0)
	if ret != nvml.SUCCESS {
		return nil
	}

	stats := &GPUStats{
		Vendor:    "NVIDIA",
		Timestamp: time.Now(),
	}

	// Get temperature
	if temp, ret := device.GetTemperature(nvml.TEMPERATURE_GPU); ret == nvml.SUCCESS {
		stats.Temperature = temp

		// Check for thermal throttling (usually starts around 83°C for most GPUs)
		if temp >= 83 {
			stats.ThrottleStatus = true
		}
	}

	// Get clock speeds
	if memoryClock, ret := device.GetClockInfo(nvml.CLOCK_MEM); ret == nvml.SUCCESS {
		stats.MemoryClock = memoryClock
	}
	if graphicsClock, ret := device.GetClockInfo(nvml.CLOCK_GRAPHICS); ret == nvml.SUCCESS {
		stats.GraphicsClock = graphicsClock
	}

	// Get memory info
	if memInfo, ret := device.GetMemoryInfo(); ret == nvml.SUCCESS {
		stats.MemoryUsed = memInfo.Used
		stats.MemoryTotal = memInfo.Total
	}

	// Get utilization
	if utilization, ret := device.GetUtilizationRates(); ret == nvml.SUCCESS {
		stats.GPUUtilization = utilization.Gpu
	}

	// Get power consumption (in milliwatts, convert to watts)
	if powerDraw, ret := device.GetPowerUsage(); ret == nvml.SUCCESS {
		stats.PowerUsage = float64(powerDraw) / 1000.0
	}

	// Get fan speed
	if fanSpeed, ret := device.GetFanSpeed(); ret == nvml.SUCCESS {
		stats.FanSpeed = fanSpeed // This is percentage, not RPM
	}

	// Check for performance state throttling
	if perfState, ret := device.GetPerformanceState(); ret == nvml.SUCCESS {
		// P0 is maximum performance, higher numbers indicate throttling
		// P2 and above usually indicate some form of throttling
		if int(perfState) > 2 {
			stats.ThrottleStatus = true
		}
	}

	return stats
}

func (app *BenchmarkApp) getGenericGPUStats() *GPUStats {
	// Try to read from common Linux GPU monitoring locations
	stats := &GPUStats{
		Timestamp: time.Now(),
	}

	// Try multiple hwmon locations for temperature
	tempLocations := []string{
		"/sys/class/hwmon/hwmon0/temp1_input",
		"/sys/class/hwmon/hwmon1/temp1_input",
		"/sys/class/hwmon/hwmon2/temp1_input",
		"/sys/class/drm/card0/device/hwmon/hwmon0/temp1_input",
		"/sys/class/drm/card0/device/hwmon/hwmon1/temp1_input",
	}

	for _, location := range tempLocations {
		if temp := app.readIntFromFile(location); temp > 0 {
			stats.Temperature = uint32(temp / 1000) // Convert from millidegrees

			// Try to determine vendor based on path
			if strings.Contains(location, "drm/card0") {
				stats.Vendor = "AMD/Intel GPU"
			} else {
				stats.Vendor = "Generic GPU"
			}

			// Check for thermal throttling
			if stats.Temperature >= 90 {
				stats.ThrottleStatus = true
			}
			break
		}
	}

	// Try Intel GPU specific location
	if stats.Temperature == 0 {
		if temp := app.readIntFromFile("/sys/class/thermal/thermal_zone0/temp"); temp > 0 {
			stats.Temperature = uint32(temp / 1000)
			stats.Vendor = "Intel GPU"
		}
	}

	// Try to read GPU power consumption (AMD specific paths)
	powerLocations := []string{
		"/sys/class/hwmon/hwmon0/power1_average",
		"/sys/class/hwmon/hwmon1/power1_average",
		"/sys/class/drm/card0/device/hwmon/hwmon0/power1_average",
		"/sys/class/drm/card0/device/hwmon/hwmon1/power1_average",
	}

	for _, location := range powerLocations {
		if power := app.readIntFromFile(location); power > 0 {
			stats.PowerUsage = float64(power) / 1000000.0 // Convert from microwatts to watts
			break
		}
	}

	// Try to read fan speed (PWM or RPM)
	fanLocations := []string{
		"/sys/class/hwmon/hwmon0/fan1_input",
		"/sys/class/hwmon/hwmon1/fan1_input",
		"/sys/class/drm/card0/device/hwmon/hwmon0/fan1_input",
	}

	for _, location := range fanLocations {
		if fanRPM := app.readIntFromFile(location); fanRPM > 0 {
			stats.FanSpeed = uint32(fanRPM)
			break
		}
	}

	// Try to read GPU clock frequencies (AMD specific)
	clockLocations := []string{
		"/sys/class/drm/card0/device/pp_dpm_sclk",
		"/sys/class/drm/card0/device/pp_dpm_mclk",
	}

	// Read GPU core clock
	if clockData := app.readStringFromFile(clockLocations[0]); clockData != "" {
		if coreClock := app.parseAMDClockInfo(clockData); coreClock > 0 {
			stats.GraphicsClock = coreClock
		}
	}

	// Read memory clock
	if clockData := app.readStringFromFile(clockLocations[1]); clockData != "" {
		if memClock := app.parseAMDClockInfo(clockData); memClock > 0 {
			stats.MemoryClock = memClock
		}
	}

	// Try to get GPU memory usage (very rough estimation)
	if memInfo := app.readMemoryInfo(); memInfo != nil {
		// This is a very rough approximation
		estimatedGPUMem := memInfo["MemTotal"] / 8 // Assume discrete GPU has 1/8 of system memory
		stats.MemoryTotal = estimatedGPUMem * 1024 // Convert to bytes

		// Estimate usage based on system memory pressure
		if memAvailable, ok := memInfo["MemAvailable"]; ok {
			memUsedSystem := memInfo["MemTotal"] - memAvailable
			usageRatio := float64(memUsedSystem) / float64(memInfo["MemTotal"])
			stats.MemoryUsed = uint64(float64(stats.MemoryTotal) * usageRatio * 0.5) // Rough estimate
		}
	}

	// If we found any meaningful data, return the stats
	if stats.Temperature > 0 || stats.PowerUsage > 0 || stats.GraphicsClock > 0 {
		if stats.Vendor == "" {
			stats.Vendor = "Generic GPU"
		}
		return stats
	}

	return nil
}

func (app *BenchmarkApp) readStringFromFile(filename string) string {
	data, err := os.ReadFile(filename)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}

func (app *BenchmarkApp) parseAMDClockInfo(clockData string) uint32 {
	// AMD clock info format: "0: 300Mhz *\n1: 600Mhz\n2: 900Mhz"
	// We want to find the active clock (marked with *)
	lines := strings.Split(clockData, "\n")
	for _, line := range lines {
		if strings.Contains(line, "*") {
			// Extract MHz value
			parts := strings.Fields(line)
			for _, part := range parts {
				if strings.HasSuffix(part, "Mhz") || strings.HasSuffix(part, "MHz") {
					clockStr := strings.TrimSuffix(strings.TrimSuffix(part, "Mhz"), "MHz")
					if clock, err := strconv.ParseUint(clockStr, 10, 32); err == nil {
						return uint32(clock)
					}
				}
			}
		}
	}
	return 0
}

func (app *BenchmarkApp) readIntFromFile(filename string) int64 {
	data, err := os.ReadFile(filename)
	if err != nil {
		return 0
	}

	value, err := strconv.ParseInt(strings.TrimSpace(string(data)), 10, 64)
	if err != nil {
		return 0
	}

	return value
}

func (app *BenchmarkApp) readMemoryInfo() map[string]uint64 {
	data, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return nil
	}

	memInfo := make(map[string]uint64)
	lines := strings.Split(string(data), "\n")

	for _, line := range lines {
		if strings.Contains(line, "MemTotal:") || strings.Contains(line, "MemAvailable:") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				key := strings.TrimSuffix(fields[0], ":")
				if value, err := strconv.ParseUint(fields[1], 10, 64); err == nil {
					memInfo[key] = value * 1024 // Convert KB to bytes
				}
			}
		}
	}

	return memInfo
}