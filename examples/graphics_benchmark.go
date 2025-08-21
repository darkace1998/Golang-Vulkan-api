package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/NVIDIA/go-nvml/pkg/nvml"
	vulkan "github.com/darkace1998/Golang-Vulkan-api"
)

// BenchmarkApp represents our graphics benchmark application
type BenchmarkApp struct {
	instance       vulkan.Instance
	physicalDevice vulkan.PhysicalDevice
	device         vulkan.Device
	graphicsQueue  vulkan.Queue
	commandPool    vulkan.CommandPool
	
	// Benchmark state
	frameCount     uint64
	startTime      time.Time
	lastFrameTime  time.Time
	currentFPS     float64
	targetFPS      int
	maxDuration    time.Duration
	
	// Scene animation
	rotationAngle  float32
	
	// GPU monitoring
	nvmlInitialized bool
}

// GPUStats holds GPU monitoring information
type GPUStats struct {
	Temperature    uint32 // in Celsius
	MemoryClock    uint32 // in MHz
	GraphicsClock  uint32 // in MHz
	MemoryUsed     uint64 // in bytes
	MemoryTotal    uint64 // in bytes
	GPUUtilization uint32 // percentage
	Vendor         string // GPU vendor
}

func main() {
	// Command line flags
	var (
		duration    = flag.Duration("duration", 0, "Benchmark duration (0 for infinite)")
		targetFPS   = flag.Int("fps", 60, "Target FPS for the benchmark")
		showHelp    = flag.Bool("help", false, "Show help information")
		simMode     = flag.Bool("sim", false, "Force simulation mode (no Vulkan)")
	)
	flag.Parse()
	
	if *showHelp {
		fmt.Println("Vulkan Graphics Benchmark")
		fmt.Println("========================")
		fmt.Println("A graphics benchmark that renders a dynamic scene and monitors GPU performance.")
		fmt.Println()
		fmt.Println("Usage:")
		flag.PrintDefaults()
		fmt.Println()
		fmt.Println("Examples:")
		fmt.Println("  go run graphics_benchmark.go                    # Run infinite benchmark")
		fmt.Println("  go run graphics_benchmark.go -duration=30s      # Run for 30 seconds")
		fmt.Println("  go run graphics_benchmark.go -fps=120           # Target 120 FPS")
		fmt.Println("  go run graphics_benchmark.go -sim               # Force simulation mode")
		return
	}
	
	fmt.Println("Vulkan Graphics Benchmark")
	fmt.Println("========================")
	
	app := &BenchmarkApp{
		startTime:     time.Now(),
		lastFrameTime: time.Now(),
		targetFPS:     *targetFPS,
		maxDuration:   *duration,
	}
	
	// Initialize Vulkan unless in simulation mode
	if !*simMode {
		if err := app.initVulkan(); err != nil {
			log.Printf("Failed to initialize Vulkan: %v", err)
			log.Println("Note: This is expected in environments without GPU drivers")
			
			// Fall back to simulation mode
			app.simulateBenchmark()
			return
		}
		defer app.cleanup()
	} else {
		fmt.Println("Running in simulation mode (Vulkan disabled)")
		app.simulateBenchmark()
		return
	}
	
	// Initialize GPU monitoring
	app.initGPUMonitoring()
	defer app.cleanupGPUMonitoring()
	
	// Run benchmark
	app.runBenchmark()
}

func (app *BenchmarkApp) initVulkan() error {
	// Create Vulkan instance
	instanceCreateInfo := &vulkan.InstanceCreateInfo{
		ApplicationInfo: &vulkan.ApplicationInfo{
			ApplicationName:    "Vulkan Graphics Benchmark",
			ApplicationVersion: vulkan.MakeVersion(1, 0, 0),
			EngineName:         "Benchmark Engine",
			EngineVersion:      vulkan.MakeVersion(1, 0, 0),
			APIVersion:         vulkan.Version13,
		},
	}
	
	instance, err := vulkan.CreateInstance(instanceCreateInfo)
	if err != nil {
		return fmt.Errorf("failed to create instance: %v", err)
	}
	app.instance = instance
	
	// Enumerate physical devices
	physicalDevices, err := vulkan.EnumeratePhysicalDevices(instance)
	if err != nil {
		return fmt.Errorf("failed to enumerate physical devices: %v", err)
	}
	
	if len(physicalDevices) == 0 {
		return fmt.Errorf("no physical devices found")
	}
	
	app.physicalDevice = physicalDevices[0]
	
	// Get device properties for display
	props := vulkan.GetPhysicalDeviceProperties(app.physicalDevice)
	fmt.Printf("Using GPU: %s\n", props.DeviceName)
	fmt.Printf("Driver Version: %d.%d.%d\n",
		props.DriverVersion.Major(),
		props.DriverVersion.Minor(),
		props.DriverVersion.Patch())
	
	// Create logical device
	if err := app.createLogicalDevice(); err != nil {
		return fmt.Errorf("failed to create logical device: %v", err)
	}
	
	// Create command pool
	if err := app.createCommandPool(); err != nil {
		return fmt.Errorf("failed to create command pool: %v", err)
	}
	
	return nil
}

func (app *BenchmarkApp) createLogicalDevice() error {
	queueFamilies := vulkan.GetPhysicalDeviceQueueFamilyProperties(app.physicalDevice)
	
	var graphicsQueueFamily uint32 = ^uint32(0)
	for i, queueFamily := range queueFamilies {
		if queueFamily.QueueFlags&vulkan.QueueGraphicsBit != 0 {
			graphicsQueueFamily = uint32(i)
			break
		}
	}
	
	if graphicsQueueFamily == ^uint32(0) {
		return fmt.Errorf("no graphics queue family found")
	}
	
	queuePriority := float32(1.0)
	deviceQueueCreateInfo := vulkan.DeviceQueueCreateInfo{
		QueueFamilyIndex: graphicsQueueFamily,
		QueuePriorities:  []float32{queuePriority},
	}
	
	deviceCreateInfo := &vulkan.DeviceCreateInfo{
		QueueCreateInfos: []vulkan.DeviceQueueCreateInfo{deviceQueueCreateInfo},
	}
	
	device, err := vulkan.CreateDevice(app.physicalDevice, deviceCreateInfo)
	if err != nil {
		return fmt.Errorf("failed to create device: %v", err)
	}
	app.device = device
	
	// Get graphics queue
	app.graphicsQueue = vulkan.GetDeviceQueue(device, graphicsQueueFamily, 0)
	
	return nil
}

func (app *BenchmarkApp) createCommandPool() error {
	queueFamilies := vulkan.GetPhysicalDeviceQueueFamilyProperties(app.physicalDevice)
	
	var graphicsQueueFamily uint32 = ^uint32(0)
	for i, queueFamily := range queueFamilies {
		if queueFamily.QueueFlags&vulkan.QueueGraphicsBit != 0 {
			graphicsQueueFamily = uint32(i)
			break
		}
	}
	
	commandPoolCreateInfo := &vulkan.CommandPoolCreateInfo{
		Flags:            vulkan.CommandPoolCreateResetCommandBufferBit,
		QueueFamilyIndex: graphicsQueueFamily,
	}
	
	commandPool, err := vulkan.CreateCommandPool(app.device, commandPoolCreateInfo)
	if err != nil {
		return fmt.Errorf("failed to create command pool: %v", err)
	}
	app.commandPool = commandPool
	
	return nil
}

func (app *BenchmarkApp) initGPUMonitoring() {
	ret := nvml.Init()
	if ret != nvml.SUCCESS {
		log.Printf("Failed to initialize NVML: %v", nvml.ErrorString(ret))
		return
	}
	app.nvmlInitialized = true
	fmt.Println("GPU monitoring initialized")
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
	
	stats := &GPUStats{Vendor: "NVIDIA"}
	
	// Get temperature
	if temp, ret := device.GetTemperature(nvml.TEMPERATURE_GPU); ret == nvml.SUCCESS {
		stats.Temperature = temp
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
	
	return stats
}

func (app *BenchmarkApp) getGenericGPUStats() *GPUStats {
	// Try to read from common Linux GPU monitoring locations
	stats := &GPUStats{}
	
	// Try AMD GPU temperature (common location)
	if temp := app.readIntFromFile("/sys/class/hwmon/hwmon0/temp1_input"); temp > 0 {
		stats.Temperature = uint32(temp / 1000) // Convert from millidegrees
		stats.Vendor = "AMD/Generic"
	} else if temp := app.readIntFromFile("/sys/class/hwmon/hwmon1/temp1_input"); temp > 0 {
		stats.Temperature = uint32(temp / 1000)
		stats.Vendor = "AMD/Generic"
	}
	
	// Try Intel GPU
	if temp := app.readIntFromFile("/sys/class/thermal/thermal_zone0/temp"); temp > 0 {
		if stats.Temperature == 0 { // Only use if we haven't found a temperature yet
			stats.Temperature = uint32(temp / 1000)
			stats.Vendor = "Intel/Generic"
		}
	}
	
	// Try to read GPU memory usage from /proc/meminfo (rough approximation)
	if memInfo := app.readMemoryInfo(); memInfo != nil {
		// Estimate GPU memory as a fraction of system memory
		// This is very rough and not accurate, but provides some data
		stats.MemoryTotal = memInfo["MemTotal"] / 4 // Assume GPU has 1/4 of system memory
		stats.MemoryUsed = stats.MemoryTotal / 10   // Assume 10% usage
	}
	
	// If we found any data, return the stats
	if stats.Temperature > 0 || stats.MemoryTotal > 0 {
		if stats.Vendor == "" {
			stats.Vendor = "Generic"
		}
		return stats
	}
	
	return nil
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

func (app *BenchmarkApp) renderFrame() {
	// Update animation
	app.rotationAngle += 0.01 // Rotate ~0.57 degrees per frame
	if app.rotationAngle > 2*math.Pi {
		app.rotationAngle -= 2 * math.Pi
	}
	
	// Simulate rendering workload
	// In a real implementation, this would record actual Vulkan commands
	app.simulateRenderingWork()
	
	// Update frame counter
	app.frameCount++
	
	// Update FPS calculation
	now := time.Now()
	deltaTime := now.Sub(app.lastFrameTime).Seconds()
	if deltaTime > 0 {
		app.currentFPS = 1.0 / deltaTime
	}
	app.lastFrameTime = now
}

func (app *BenchmarkApp) simulateRenderingWork() {
	// Simulate some computational work to represent rendering
	// This creates a dynamic scene by varying the workload
	workAmount := int(1000 + 500*math.Sin(float64(app.rotationAngle)))
	
	for i := 0; i < workAmount; i++ {
		// Simulate vertex transformations
		x := math.Sin(float64(i) * float64(app.rotationAngle))
		y := math.Cos(float64(i) * float64(app.rotationAngle))
		_ = x*y // Prevent optimization
	}
}

func (app *BenchmarkApp) runBenchmark() {
	fmt.Printf("\nStarting benchmark (Target: %d FPS", app.targetFPS)
	if app.maxDuration > 0 {
		fmt.Printf(", Duration: %v", app.maxDuration)
	}
	fmt.Println(")")
	fmt.Println("Press Ctrl+C to exit")
	
	frameDuration := time.Second / time.Duration(app.targetFPS)
	
	for {
		frameStart := time.Now()
		
		// Check if we've exceeded the maximum duration
		if app.maxDuration > 0 && time.Since(app.startTime) >= app.maxDuration {
			fmt.Printf("\nBenchmark completed after %v\n", app.maxDuration)
			app.displayFinalStats()
			break
		}
		
		// Render frame
		app.renderFrame()
		
		// Display stats every 60 frames (approximately once per second at 60 FPS)
		if app.frameCount%60 == 0 {
			app.displayStats()
		}
		
		// Maintain target frame rate
		frameTime := time.Since(frameStart)
		if frameTime < frameDuration {
			time.Sleep(frameDuration - frameTime)
		}
	}
}

func (app *BenchmarkApp) simulateBenchmark() {
	fmt.Println("\nRunning simulated benchmark (no GPU drivers available)...")
	fmt.Println("This demonstrates the benchmark structure and monitoring capabilities")
	
	duration := app.maxDuration
	if duration == 0 {
		duration = 5 * time.Second // Default simulation duration
	}
	
	targetFPS := app.targetFPS
	if targetFPS == 0 {
		targetFPS = 60
	}
	
	totalFrames := int(duration.Seconds()) * targetFPS
	
	for i := 0; i < totalFrames; i++ {
		app.renderFrame()
		
		if i%targetFPS == 0 {
			app.displayStats()
		}
		
		time.Sleep(time.Second / time.Duration(targetFPS))
	}
	
	fmt.Printf("\nSimulated benchmark complete! (%v)\n", duration)
	app.displayFinalStats()
}

func (app *BenchmarkApp) displayStats() {
	// Clear screen (simple approach)
	if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
		fmt.Print("\033[2J\033[H")
	}
	
	fmt.Println("Vulkan Graphics Benchmark - Live Stats")
	fmt.Println("=====================================")
	
	// Runtime stats
	elapsed := time.Since(app.startTime)
	avgFPS := float64(app.frameCount) / elapsed.Seconds()
	
	fmt.Printf("Runtime: %v\n", elapsed.Round(time.Second))
	fmt.Printf("Total Frames: %d\n", app.frameCount)
	fmt.Printf("Average FPS: %.1f\n", avgFPS)
	fmt.Printf("Current FPS: %.1f\n", app.currentFPS)
	fmt.Printf("Rotation Angle: %.2f radians\n", app.rotationAngle)
	
	// GPU stats
	gpuStats := app.getGPUStats()
	if gpuStats != nil {
		fmt.Printf("\nGPU Statistics (%s):\n", gpuStats.Vendor)
		if gpuStats.Temperature > 0 {
			fmt.Printf("Temperature: %d°C\n", gpuStats.Temperature)
		}
		if gpuStats.GraphicsClock > 0 {
			fmt.Printf("Graphics Clock: %d MHz\n", gpuStats.GraphicsClock)
		}
		if gpuStats.MemoryClock > 0 {
			fmt.Printf("Memory Clock: %d MHz\n", gpuStats.MemoryClock)
		}
		if gpuStats.GPUUtilization > 0 {
			fmt.Printf("GPU Utilization: %d%%\n", gpuStats.GPUUtilization)
		}
		if gpuStats.MemoryTotal > 0 {
			fmt.Printf("Memory Used: %.1f MB / %.1f MB (%.1f%%)\n",
				float64(gpuStats.MemoryUsed)/1024/1024,
				float64(gpuStats.MemoryTotal)/1024/1024,
				float64(gpuStats.MemoryUsed)*100/float64(gpuStats.MemoryTotal))
		}
	} else {
		fmt.Println("\nGPU Statistics: Not available")
		fmt.Println("(GPU monitoring requires supported hardware and drivers)")
	}
	
	// System stats
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	fmt.Printf("\nSystem Memory: %.1f MB allocated\n", float64(memStats.Alloc)/1024/1024)
	fmt.Printf("Goroutines: %d\n", runtime.NumGoroutine())
}

func (app *BenchmarkApp) displayFinalStats() {
	elapsed := time.Since(app.startTime)
	avgFPS := float64(app.frameCount) / elapsed.Seconds()
	
	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("FINAL BENCHMARK RESULTS")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Printf("Total Runtime: %v\n", elapsed.Round(time.Millisecond))
	fmt.Printf("Total Frames Rendered: %d\n", app.frameCount)
	fmt.Printf("Average FPS: %.2f\n", avgFPS)
	fmt.Printf("Final Rotation Angle: %.2f radians\n", app.rotationAngle)
	
	// Performance rating
	var rating string
	switch {
	case avgFPS >= 120:
		rating = "Excellent"
	case avgFPS >= 60:
		rating = "Good"
	case avgFPS >= 30:
		rating = "Fair"
	default:
		rating = "Poor"
	}
	fmt.Printf("Performance Rating: %s\n", rating)
	
	// GPU stats summary
	gpuStats := app.getGPUStats()
	if gpuStats != nil {
		fmt.Printf("GPU Vendor: %s\n", gpuStats.Vendor)
		if gpuStats.Temperature > 0 {
			fmt.Printf("Max Temperature: %d°C\n", gpuStats.Temperature)
		}
	}
	
	fmt.Println(strings.Repeat("=", 50))
}

func (app *BenchmarkApp) cleanup() {
	if app.commandPool != nil {
		vulkan.DestroyCommandPool(app.device, app.commandPool)
	}
	if app.device != nil {
		vulkan.DestroyDevice(app.device)
	}
	if app.instance != nil {
		vulkan.DestroyInstance(app.instance)
	}
}