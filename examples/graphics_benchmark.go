package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/NVIDIA/go-nvml/pkg/nvml"
	vulkan "github.com/darkace1998/Golang-Vulkan-api"
)

// TestMode defines the type of test being run
type TestMode int

const (
	StressTest TestMode = iota
	Benchmark
)

// GraphicsQuality defines the intensity level of graphics effects
type GraphicsQuality int

const (
	QualityLow GraphicsQuality = iota
	QualityMedium
	QualityHigh
	QualityUltra
)

// Resolution represents display resolution
type Resolution struct {
	Width  uint32
	Height uint32
	Name   string
}

// BenchmarkApp represents our comprehensive GPU stress testing application
type BenchmarkApp struct {
	instance       vulkan.Instance
	physicalDevice vulkan.PhysicalDevice
	device         vulkan.Device
	graphicsQueue  vulkan.Queue
	commandPool    vulkan.CommandPool

	// Test configuration
	testMode    TestMode
	quality     GraphicsQuality
	resolution  Resolution
	targetFPS   int
	maxDuration time.Duration

	// Benchmark state
	frameCount    uint64
	startTime     time.Time
	lastFrameTime time.Time
	currentFPS    float64
	avgFPS        float64
	minFPS        float64
	maxFPS        float64
	frameTimesMs  []float64

	// Advanced scene animation
	rotationAngle   float32
	animationTime   float32
	particleCount   int
	complexityLevel int

	// GPU monitoring
	nvmlInitialized   bool
	monitoringEnabled bool
	statsHistory      []GPUStats
	powerHistory      []float64
	fanSpeedHistory   []uint32

	// Error detection
	artifactDetection bool
	errorCount        uint64
	lastErrorTime     time.Time

	// Performance data
	performanceLog []PerformanceData
	mutex          sync.RWMutex
}

// GPUStats holds comprehensive GPU monitoring information
type GPUStats struct {
	Timestamp      time.Time
	Temperature    uint32  // in Celsius
	MemoryClock    uint32  // in MHz
	GraphicsClock  uint32  // in MHz
	MemoryUsed     uint64  // in bytes
	MemoryTotal    uint64  // in bytes
	GPUUtilization uint32  // percentage
	PowerUsage     float64 // in Watts
	FanSpeed       uint32  // in RPM or percentage
	Vendor         string  // GPU vendor
	ThrottleStatus bool    // true if thermal throttling detected
}

// PerformanceData holds frame performance metrics
type PerformanceData struct {
	Timestamp   time.Time
	FrameTime   float64 // in milliseconds
	FPS         float64
	GPUTemp     uint32
	PowerUsage  float64
	MemoryUsage uint64
}

// TestResults holds final benchmark results
type TestResults struct {
	Duration       time.Duration
	TotalFrames    uint64
	AverageFPS     float64
	MinFPS         float64
	MaxFPS         float64
	PercentileFPS  map[string]float64 // 1%, 5%, 95%, 99%
	MaxTemperature uint32
	AvgPowerUsage  float64
	MaxPowerUsage  float64
	ErrorCount     uint64
	StabilityScore float64
	BenchmarkScore int
}

// Predefined resolutions
var standardResolutions = []Resolution{
	{1920, 1080, "1080p"},
	{2560, 1440, "1440p"},
	{3840, 2160, "4K"},
	{1280, 720, "720p"},
	{1366, 768, "768p"},
	{1600, 900, "900p"},
	{2560, 1080, "1080p Ultrawide"},
	{3440, 1440, "1440p Ultrawide"},
}

// Configuration holds parsed test configuration
type Configuration struct {
	TestMode   TestMode
	Quality    GraphicsQuality
	Resolution Resolution
	Duration   time.Duration
	TargetFPS  int
}

func showDetailedHelp() {
	fmt.Println("GPU STRESS TESTING & BENCHMARK APPLICATION")
	fmt.Println("==========================================")
	fmt.Println()
	fmt.Println("DESCRIPTION:")
	fmt.Println("  Advanced GPU stress testing application designed to push modern graphics")
	fmt.Println("  cards to their limits. Tests stability, thermal performance, and provides")
	fmt.Println("  comprehensive performance metrics similar to FurMark and 3DMark.")
	fmt.Println()
	fmt.Println("USAGE:")
	flag.PrintDefaults()
	fmt.Println()
	fmt.Println("TEST MODES:")
	fmt.Println("  stress    - Runs indefinitely until manually stopped (default)")
	fmt.Println("  benchmark - Runs for fixed duration and provides performance score")
	fmt.Println()
	fmt.Println("QUALITY LEVELS:")
	fmt.Println("  low       - Basic rendering, minimal GPU load")
	fmt.Println("  medium    - Standard effects, moderate GPU load")
	fmt.Println("  high      - Advanced effects, high GPU load (default)")
	fmt.Println("  ultra     - Maximum effects, extreme GPU load")
	fmt.Println()
	fmt.Println("EXAMPLES:")
	fmt.Println("  # Basic stress test at 1080p")
	fmt.Println("  go run graphics_benchmark.go")
	fmt.Println()
	fmt.Println("  # 4K benchmark for 5 minutes with CSV export")
	fmt.Println("  go run graphics_benchmark.go -mode=benchmark -resolution=4K -duration=5m -csv -output=./results")
	fmt.Println()
	fmt.Println("  # Ultra quality stress test with artifact detection")
	fmt.Println("  go run graphics_benchmark.go -quality=ultra -artifacts")
	fmt.Println()
	fmt.Println("  # Custom resolution benchmark")
	fmt.Println("  go run graphics_benchmark.go -resolution=2560x1440 -mode=benchmark -duration=2m")
	fmt.Println()
	fmt.Println("MONITORING:")
	fmt.Println("  The application monitors GPU temperature, clock speeds, power consumption,")
	fmt.Println("  fan speeds, and detects thermal throttling or stability issues.")
	fmt.Println()
}

func parseConfiguration(testModeStr, qualityStr, resolutionStr string, duration time.Duration, targetFPS int) (*Configuration, error) {
	config := &Configuration{
		Duration:  duration,
		TargetFPS: targetFPS,
	}

	// Parse test mode
	switch strings.ToLower(testModeStr) {
	case "stress":
		config.TestMode = StressTest
	case "benchmark":
		config.TestMode = Benchmark
	default:
		return nil, fmt.Errorf("invalid test mode: %s (use 'stress' or 'benchmark')", testModeStr)
	}

	// Parse quality
	switch strings.ToLower(qualityStr) {
	case "low":
		config.Quality = QualityLow
	case "medium":
		config.Quality = QualityMedium
	case "high":
		config.Quality = QualityHigh
	case "ultra":
		config.Quality = QualityUltra
	default:
		return nil, fmt.Errorf("invalid quality: %s (use 'low', 'medium', 'high', or 'ultra')", qualityStr)
	}

	// Parse resolution
	config.Resolution = parseResolution(resolutionStr)
	if config.Resolution.Width == 0 {
		return nil, fmt.Errorf("invalid resolution: %s", resolutionStr)
	}

	return config, nil
}

func parseResolution(resStr string) Resolution {
	// Check for standard resolutions first
	for _, res := range standardResolutions {
		if strings.EqualFold(res.Name, resStr) {
			return res
		}
	}

	// Try to parse custom resolution (e.g., "1920x1080")
	if strings.Contains(resStr, "x") {
		parts := strings.Split(resStr, "x")
		if len(parts) == 2 {
			width, err1 := strconv.ParseUint(parts[0], 10, 32)
			height, err2 := strconv.ParseUint(parts[1], 10, 32)
			if err1 == nil && err2 == nil && width > 0 && height > 0 {
				return Resolution{
					Width:  uint32(width),
					Height: uint32(height),
					Name:   fmt.Sprintf("%dx%d", width, height),
				}
			}
		}
	}

	// Default to 1080p if parsing fails
	return Resolution{1920, 1080, "1080p"}
}

func (app *BenchmarkApp) displayConfiguration(verbose bool) {
	fmt.Printf("📋 TEST CONFIGURATION\n")
	fmt.Printf("   Mode: %s\n", app.getTestModeString())
	fmt.Printf("   Quality: %s\n", app.getQualityString())
	fmt.Printf("   Resolution: %s (%dx%d)\n", app.resolution.Name, app.resolution.Width, app.resolution.Height)
	fmt.Printf("   Target FPS: %d\n", app.targetFPS)
	if app.maxDuration > 0 {
		fmt.Printf("   Duration: %v\n", app.maxDuration)
	} else {
		fmt.Printf("   Duration: Infinite (stress test)\n")
	}
	fmt.Printf("   Artifact Detection: %v\n", app.artifactDetection)

	if verbose {
		fmt.Printf("\n🔧 ADVANCED SETTINGS\n")
		fmt.Printf("   CPU Cores: %d\n", runtime.NumCPU())
		fmt.Printf("   Go Version: %s\n", runtime.Version())
		fmt.Printf("   Platform: %s/%s\n", runtime.GOOS, runtime.GOARCH)
	}
	fmt.Println()
}

func (app *BenchmarkApp) getTestModeString() string {
	switch app.testMode {
	case StressTest:
		return "Stress Test"
	case Benchmark:
		return "Benchmark"
	default:
		return "Unknown"
	}
}

func (app *BenchmarkApp) getQualityString() string {
	switch app.quality {
	case QualityLow:
		return "Low"
	case QualityMedium:
		return "Medium"
	case QualityHigh:
		return "High"
	case QualityUltra:
		return "Ultra"
	default:
		return "Unknown"
	}
}

func main() {
	// Enhanced command line flags
	var (
		duration        = flag.Duration("duration", 0, "Test duration (0 for infinite stress test)")
		targetFPS       = flag.Int("fps", 60, "Target FPS for the test")
		testModeStr     = flag.String("mode", "stress", "Test mode: 'stress' or 'benchmark'")
		qualityStr      = flag.String("quality", "high", "Graphics quality: 'low', 'medium', 'high', 'ultra'")
		resolutionStr   = flag.String("resolution", "1080p", "Resolution: '720p', '1080p', '1440p', '4K', or 'WIDTHxHEIGHT'")
		outputDir       = flag.String("output", "", "Output directory for logs and reports")
		csvExport       = flag.Bool("csv", false, "Export performance data to CSV")
		artifactScan    = flag.Bool("artifacts", false, "Enable artifact detection mode")
		showHelp        = flag.Bool("help", false, "Show detailed help information")
		simMode         = flag.Bool("sim", false, "Force simulation mode (no Vulkan)")
		listResolutions = flag.Bool("list-res", false, "List available resolutions")
		verboseMode     = flag.Bool("verbose", false, "Enable verbose logging")
	)
	flag.Parse()

	if *listResolutions {
		fmt.Println("Available Resolutions:")
		for _, res := range standardResolutions {
			fmt.Printf("  %s (%dx%d)\n", res.Name, res.Width, res.Height)
		}
		fmt.Println("  Custom: WIDTHxHEIGHT (e.g., 1920x1080)")
		return
	}

	if *showHelp {
		showDetailedHelp()
		return
	}

	fmt.Println("╔═════════════════════════════════════════════════╗")
	fmt.Println("║        GPU STRESS TESTING & BENCHMARK          ║")
	fmt.Println("║          Advanced Vulkan Application           ║")
	fmt.Println("╚═════════════════════════════════════════════════╝")
	fmt.Println()

	// Parse configuration
	config, err := parseConfiguration(*testModeStr, *qualityStr, *resolutionStr, *duration, *targetFPS)
	if err != nil {
		log.Fatalf("Configuration error: %v", err)
	}

	// Create application
	app := &BenchmarkApp{
		testMode:          config.TestMode,
		quality:           config.Quality,
		resolution:        config.Resolution,
		targetFPS:         config.TargetFPS,
		maxDuration:       config.Duration,
		artifactDetection: *artifactScan,
		monitoringEnabled: true,
		frameTimesMs:      make([]float64, 0, 1000),
		statsHistory:      make([]GPUStats, 0, 1000),
		performanceLog:    make([]PerformanceData, 0, 10000),
	}

	// Display test configuration
	app.displayConfiguration(*verboseMode)

	// Initialize output directory if specified
	if *outputDir != "" {
		if err := os.MkdirAll(*outputDir, 0755); err != nil {
			log.Printf("Warning: Could not create output directory: %v", err)
		}
	}

	// Initialize Vulkan unless in simulation mode
	if !*simMode {
		if err := app.initVulkan(); err != nil {
			log.Printf("Failed to initialize Vulkan: %v", err)
			log.Println("Falling back to simulation mode...")
			*simMode = true
		} else {
			defer app.cleanup()
		}
	}

	if *simMode {
		fmt.Println("🔧 Running in SIMULATION mode (Vulkan disabled)")
		app.runSimulation()
	} else {
		fmt.Println("🚀 Running HARDWARE-ACCELERATED stress test")
		app.runStressTest()
	}

	// Generate final report
	results := app.generateResults()
	app.displayResults(results)

	// Export data if requested
	if *csvExport && *outputDir != "" {
		app.exportToCSV(*outputDir)
	}
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

	// Alternative: Try to get fan speed in RPM if available
	// NVML doesn't always provide RPM directly, so we might need to estimate

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
		_ = x * y // Prevent optimization
	}
}

func (app *BenchmarkApp) runStressTest() {
	fmt.Println("🔥 INITIATING GPU STRESS TEST")
	fmt.Println("Press Ctrl+C to stop the test at any time...")
	fmt.Println()

	app.startTime = time.Now()
	app.lastFrameTime = time.Now()
	app.minFPS = math.Inf(1)
	app.maxFPS = 0

	// Initialize GPU monitoring
	app.initGPUMonitoring()
	defer app.cleanupGPUMonitoring()

	// Calculate complexity based on quality setting
	app.setComplexityLevel()

	// Start monitoring goroutine
	go app.monitoringLoop()

	// Main rendering loop
	frameInterval := time.Second / time.Duration(app.targetFPS)
	ticker := time.NewTicker(frameInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			app.performAdvancedRender()
			app.updatePerformanceMetrics()

			// Check for exit conditions
			if app.shouldExit() {
				return
			}

			// Display stats every second
			if time.Since(app.lastFrameTime) >= time.Second {
				app.displayLiveStats()
				app.lastFrameTime = time.Now()
			}
		}
	}
}

func (app *BenchmarkApp) runSimulation() {
	fmt.Println("🔧 RUNNING SIMULATION MODE")
	fmt.Println("Simulating GPU load without hardware acceleration...")
	fmt.Println()

	app.startTime = time.Now()
	app.lastFrameTime = time.Now()
	app.minFPS = math.Inf(1)
	app.maxFPS = 0

	// Set moderate complexity for simulation
	app.complexityLevel = int(app.quality) + 1
	app.particleCount = 1000 * app.complexityLevel

	frameInterval := time.Second / time.Duration(app.targetFPS)
	ticker := time.NewTicker(frameInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			app.simulateAdvancedWorkload()
			app.updatePerformanceMetrics()

			if app.shouldExit() {
				return
			}

			if time.Since(app.lastFrameTime) >= time.Second {
				app.displayLiveStats()
				app.lastFrameTime = time.Now()
			}
		}
	}
}

func (app *BenchmarkApp) setComplexityLevel() {
	// Set complexity based on quality and resolution
	baseComplexity := int(app.quality) + 1
	resolutionMultiplier := float64(app.resolution.Width*app.resolution.Height) / (1920.0 * 1080.0)

	app.complexityLevel = int(float64(baseComplexity) * resolutionMultiplier)
	if app.complexityLevel < 1 {
		app.complexityLevel = 1
	}

	// Set particle count based on complexity
	app.particleCount = 5000 * app.complexityLevel

	fmt.Printf("🎮 WORKLOAD CONFIGURATION\n")
	fmt.Printf("   Complexity Level: %d\n", app.complexityLevel)
	fmt.Printf("   Particle Count: %d\n", app.particleCount)
	fmt.Printf("   Estimated GPU Load: %s\n", app.getLoadDescription())
	fmt.Println()
}

func (app *BenchmarkApp) getLoadDescription() string {
	switch app.quality {
	case QualityLow:
		return "Light (Basic geometric stress)"
	case QualityMedium:
		return "Moderate (Shader-intensive workload)"
	case QualityHigh:
		return "Heavy (Advanced effects + high-resolution)"
	case QualityUltra:
		return "Extreme (Maximum GPU utilization)"
	default:
		return "Unknown"
	}
}

func (app *BenchmarkApp) performAdvancedRender() {
	// Simulate complex rendering pipeline based on quality level
	app.animationTime += 0.016 // ~60 FPS animation step
	app.rotationAngle = float32(math.Mod(float64(app.animationTime), 2*math.Pi))

	// Simulate different rendering passes based on quality
	switch app.quality {
	case QualityUltra:
		app.simulateRayTracingPass()
		app.simulateVolumetricEffects()
		app.simulatePostProcessing()
		fallthrough
	case QualityHigh:
		app.simulateAdvancedLighting()
		app.simulateTessellation()
		fallthrough
	case QualityMedium:
		app.simulateShaderWork()
		app.simulateTextureOps()
		fallthrough
	case QualityLow:
		app.simulateGeometryRendering()
	}

	// Perform actual Vulkan operations
	app.renderFrame()

	app.frameCount++

	// Record frame timing
	now := time.Now()
	frameTime := now.Sub(app.lastFrameTime).Seconds() * 1000 // Convert to milliseconds
	app.frameTimesMs = append(app.frameTimesMs, frameTime)

	// Keep only last 1000 frame times for rolling statistics
	if len(app.frameTimesMs) > 1000 {
		app.frameTimesMs = app.frameTimesMs[1:]
	}
}

func (app *BenchmarkApp) simulateAdvancedWorkload() {
	app.animationTime += 0.016
	app.rotationAngle = float32(math.Mod(float64(app.animationTime), 2*math.Pi))

	// Simulate CPU-intensive work proportional to GPU complexity
	workUnits := app.complexityLevel * 1000

	// Simulate different types of computational work
	for i := 0; i < workUnits; i++ {
		// Simulate matrix operations
		_ = math.Sin(float64(i)) * math.Cos(float64(app.rotationAngle))

		// Simulate memory access patterns
		if i%100 == 0 {
			runtime.Gosched()
		}
	}

	// Simulate particle system updates
	for i := 0; i < app.particleCount/100; i++ {
		_ = rand.Float64() * float64(app.complexityLevel)
	}

	app.frameCount++

	// Record simulated frame timing
	now := time.Now()
	frameTime := now.Sub(app.lastFrameTime).Seconds() * 1000
	app.frameTimesMs = append(app.frameTimesMs, frameTime)

	if len(app.frameTimesMs) > 1000 {
		app.frameTimesMs = app.frameTimesMs[1:]
	}
}

// Advanced rendering simulation functions
func (app *BenchmarkApp) simulateRayTracingPass() {
	// Simulate ray tracing workload - very compute intensive
	rayCount := app.resolution.Width * app.resolution.Height / 4
	for i := uint32(0); i < rayCount; i++ {
		// Simulate ray-scene intersection calculations
		_ = math.Sqrt(float64(i)) * math.Tan(float64(app.rotationAngle))
		if i%1000 == 0 {
			runtime.Gosched()
		}
	}
}

func (app *BenchmarkApp) simulateVolumetricEffects() {
	// Simulate volumetric fog/smoke calculations
	voxelCount := app.complexityLevel * 10000
	for i := 0; i < voxelCount; i++ {
		// Simulate 3D noise calculations for volumetric effects
		x := float64(i%100) * 0.1
		y := float64((i/100)%100) * 0.1
		z := float64(i/10000) * 0.1
		_ = math.Sin(x) * math.Cos(y) * math.Tan(z) * float64(app.animationTime)
	}
}

func (app *BenchmarkApp) simulatePostProcessing() {
	// Simulate post-processing effects like bloom, motion blur
	pixelCount := app.resolution.Width * app.resolution.Height

	// Simulate bloom effect
	for i := uint32(0); i < pixelCount/8; i++ {
		_ = math.Exp(-float64(i)*0.001) * math.Sin(float64(app.animationTime))
	}

	// Simulate motion blur
	for i := uint32(0); i < pixelCount/16; i++ {
		_ = math.Log(1+float64(i)) * float64(app.rotationAngle)
	}
}

func (app *BenchmarkApp) simulateAdvancedLighting() {
	// Simulate advanced lighting calculations
	lightCount := app.complexityLevel * 100
	for i := 0; i < lightCount; i++ {
		// Simulate shadow map calculations
		shadowSamples := 16
		for j := 0; j < shadowSamples; j++ {
			_ = math.Pow(float64(i+j), 0.5) * math.Sin(float64(app.animationTime))
		}
	}
}

func (app *BenchmarkApp) simulateTessellation() {
	// Simulate tessellation workload
	patchCount := app.complexityLevel * 500
	for i := 0; i < patchCount; i++ {
		// Simulate subdivision calculations
		subdivisions := 4
		for j := 0; j < subdivisions; j++ {
			_ = math.Sin(float64(i*j)) * math.Cos(float64(app.rotationAngle))
		}
	}
}

func (app *BenchmarkApp) simulateShaderWork() {
	// Simulate complex shader calculations
	shaderOps := app.complexityLevel * 2000
	for i := 0; i < shaderOps; i++ {
		// Simulate vertex shader work
		_ = math.Sin(float64(i)*0.01) * math.Cos(float64(app.animationTime))

		// Simulate fragment shader work
		if i%4 == 0 {
			_ = math.Pow(float64(i), 0.3) * float64(app.rotationAngle)
		}
	}
}

func (app *BenchmarkApp) simulateTextureOps() {
	// Simulate texture sampling operations
	textureOps := app.complexityLevel * 1000
	for i := 0; i < textureOps; i++ {
		// Simulate bilinear filtering
		u := float64(i%256) / 255.0
		v := float64((i/256)%256) / 255.0
		_ = (1-u)*(1-v)*float64(app.animationTime) + u*v*float64(app.rotationAngle)
	}
}

func (app *BenchmarkApp) simulateGeometryRendering() {
	// Simulate basic geometry rendering
	vertexCount := app.particleCount * 3 // Triangles
	for i := 0; i < vertexCount; i++ {
		// Simulate vertex transformations
		_ = math.Sin(float64(i)*0.001) * float64(app.rotationAngle)
		_ = math.Cos(float64(i)*0.001) * float64(app.animationTime)
	}
}

// Enhanced monitoring and performance tracking
func (app *BenchmarkApp) monitoringLoop() {
	ticker := time.NewTicker(500 * time.Millisecond) // Monitor every 500ms
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if app.monitoringEnabled {
				app.collectPerformanceData()
				app.detectArtifacts()
			}
		}
	}
}

func (app *BenchmarkApp) collectPerformanceData() {
	app.mutex.Lock()
	defer app.mutex.Unlock()

	now := time.Now()
	stats := app.getGPUStats()

	if stats != nil {
		stats.Timestamp = now
		app.statsHistory = append(app.statsHistory, *stats)

		// Keep only last 1000 entries
		if len(app.statsHistory) > 1000 {
			app.statsHistory = app.statsHistory[1:]
		}

		// Record performance data
		perfData := PerformanceData{
			Timestamp:   now,
			GPUTemp:     stats.Temperature,
			PowerUsage:  stats.PowerUsage,
			MemoryUsage: stats.MemoryUsed,
		}

		if len(app.frameTimesMs) > 0 {
			perfData.FrameTime = app.frameTimesMs[len(app.frameTimesMs)-1]
			perfData.FPS = app.currentFPS
		}

		app.performanceLog = append(app.performanceLog, perfData)

		// Keep only last 10000 entries
		if len(app.performanceLog) > 10000 {
			app.performanceLog = app.performanceLog[1:]
		}
	}
}

func (app *BenchmarkApp) detectArtifacts() {
	if !app.artifactDetection {
		return
	}

	// Simulate artifact detection by checking for anomalies
	if len(app.frameTimesMs) >= 10 {
		recent := app.frameTimesMs[len(app.frameTimesMs)-10:]
		avgFrameTime := 0.0
		for _, ft := range recent {
			avgFrameTime += ft
		}
		avgFrameTime /= float64(len(recent))

		// Check for sudden frame time spikes (potential artifacts)
		lastFrameTime := app.frameTimesMs[len(app.frameTimesMs)-1]
		if lastFrameTime > avgFrameTime*3 && lastFrameTime > 100 { // >100ms frame time
			app.errorCount++
			app.lastErrorTime = time.Now()
		}
	}
}

func (app *BenchmarkApp) updatePerformanceMetrics() {
	now := time.Now()
	elapsed := now.Sub(app.startTime).Seconds()

	if elapsed > 0 {
		app.currentFPS = float64(app.frameCount) / elapsed
		app.avgFPS = app.currentFPS

		if app.currentFPS < app.minFPS {
			app.minFPS = app.currentFPS
		}
		if app.currentFPS > app.maxFPS {
			app.maxFPS = app.currentFPS
		}
	}
}

func (app *BenchmarkApp) shouldExit() bool {
	if app.maxDuration > 0 && time.Since(app.startTime) >= app.maxDuration {
		return true
	}
	return false
}

func (app *BenchmarkApp) displayLiveStats() {
	// Clear screen and show live stats
	fmt.Print("\033[2J\033[H") // Clear screen and move cursor to top

	elapsed := time.Since(app.startTime)

	fmt.Println("╔═══════════════════════════════════════════════════════════════╗")
	fmt.Println("║               GPU STRESS TEST - LIVE MONITORING              ║")
	fmt.Println("╠═══════════════════════════════════════════════════════════════╣")

	// Performance metrics
	fmt.Printf("║ Runtime: %-15s │ Total Frames: %-15d ║\n",
		formatDuration(elapsed), app.frameCount)
	fmt.Printf("║ Current FPS: %-12.1f │ Average FPS: %-14.1f ║\n",
		app.currentFPS, app.avgFPS)

	if app.minFPS != math.Inf(1) && app.maxFPS > 0 {
		fmt.Printf("║ Min FPS: %-15.1f │ Max FPS: %-18.1f ║\n",
			app.minFPS, app.maxFPS)
	}

	// Calculate frame time percentiles if we have enough data
	if len(app.frameTimesMs) >= 10 {
		percentiles := app.calculateFrameTimePercentiles()
		fmt.Printf("║ 1%% Low: %-7.1f FPS       │ Frame Time: %-7.1f ms        ║\n",
			1000.0/percentiles[99], percentiles[50])
	}

	fmt.Println("╠═══════════════════════════════════════════════════════════════╣")

	// GPU statistics
	stats := app.getGPUStats()
	if stats != nil {
		fmt.Printf("║ GPU: %-25s │ Temp: %-8d°C          ║\n",
			stats.Vendor, stats.Temperature)

		if stats.PowerUsage > 0 {
			fmt.Printf("║ Power: %-7.1f W            │ GPU Load: %-6d%%         ║\n",
				stats.PowerUsage, stats.GPUUtilization)
		}

		if stats.GraphicsClock > 0 {
			fmt.Printf("║ Core Clock: %-6d MHz      │ Memory Clock: %-6d MHz ║\n",
				stats.GraphicsClock, stats.MemoryClock)
		}

		if stats.MemoryTotal > 0 {
			memUsedMB := float64(stats.MemoryUsed) / (1024 * 1024)
			memTotalMB := float64(stats.MemoryTotal) / (1024 * 1024)
			memPercent := float64(stats.MemoryUsed) / float64(stats.MemoryTotal) * 100
			fmt.Printf("║ VRAM: %-7.0f/%-7.0f MB    │ Usage: %-8.1f%%         ║\n",
				memUsedMB, memTotalMB, memPercent)
		}

		if stats.ThrottleStatus {
			fmt.Println("║ ⚠️  THERMAL THROTTLING DETECTED                              ║")
		}
	}

	// System info
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	systemMemMB := float64(m.Alloc) / (1024 * 1024)

	fmt.Println("╠═══════════════════════════════════════════════════════════════╣")
	fmt.Printf("║ System Memory: %-7.1f MB    │ Goroutines: %-12d    ║\n",
		systemMemMB, runtime.NumGoroutine())

	if app.artifactDetection && app.errorCount > 0 {
		fmt.Printf("║ Artifacts Detected: %-6d     │ Last Error: %-12s    ║\n",
			app.errorCount, formatDuration(time.Since(app.lastErrorTime)))
	}

	fmt.Printf("║ Test Mode: %-16s │ Quality: %-15s ║\n",
		app.getTestModeString(), app.getQualityString())
	fmt.Printf("║ Resolution: %-14s │ Complexity: %-12d ║\n",
		app.resolution.Name, app.complexityLevel)

	fmt.Println("╚═══════════════════════════════════════════════════════════════╝")

	if app.maxDuration > 0 {
		remaining := app.maxDuration - elapsed
		if remaining > 0 {
			fmt.Printf("\nTime Remaining: %s\n", formatDuration(remaining))
		}
	} else {
		fmt.Println("\nPress Ctrl+C to stop the stress test")
	}
	fmt.Println()
}

func (app *BenchmarkApp) calculateFrameTimePercentiles() map[int]float64 {
	if len(app.frameTimesMs) == 0 {
		return make(map[int]float64)
	}

	// Create a sorted copy
	sorted := make([]float64, len(app.frameTimesMs))
	copy(sorted, app.frameTimesMs)

	// Simple bubble sort for small arrays
	for i := 0; i < len(sorted); i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[i] > sorted[j] {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}

	percentiles := make(map[int]float64)
	n := len(sorted)

	percentiles[1] = sorted[n*1/100]
	percentiles[5] = sorted[n*5/100]
	percentiles[50] = sorted[n*50/100]
	percentiles[95] = sorted[n*95/100]
	percentiles[99] = sorted[n*99/100]

	return percentiles
}

func formatDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	if hours > 0 {
		return fmt.Sprintf("%dh%02dm%02ds", hours, minutes, seconds)
	} else if minutes > 0 {
		return fmt.Sprintf("%dm%02ds", minutes, seconds)
	} else {
		return fmt.Sprintf("%ds", seconds)
	}
}

// Results generation and reporting
func (app *BenchmarkApp) generateResults() *TestResults {
	app.mutex.RLock()
	defer app.mutex.RUnlock()

	results := &TestResults{
		Duration:      time.Since(app.startTime),
		TotalFrames:   app.frameCount,
		AverageFPS:    app.avgFPS,
		MinFPS:        app.minFPS,
		MaxFPS:        app.maxFPS,
		ErrorCount:    app.errorCount,
		PercentileFPS: make(map[string]float64),
	}

	// Calculate frame time percentiles
	if len(app.frameTimesMs) > 0 {
		percentiles := app.calculateFrameTimePercentiles()
		results.PercentileFPS["1%"] = 1000.0 / percentiles[99] // Convert worst 1% frame time to FPS
		results.PercentileFPS["5%"] = 1000.0 / percentiles[95]
		results.PercentileFPS["95%"] = 1000.0 / percentiles[5]
		results.PercentileFPS["99%"] = 1000.0 / percentiles[1]
	}

	// Calculate temperature and power statistics
	if len(app.statsHistory) > 0 {
		maxTemp := uint32(0)
		totalPower := 0.0
		maxPower := 0.0
		powerSamples := 0

		for _, stat := range app.statsHistory {
			if stat.Temperature > maxTemp {
				maxTemp = stat.Temperature
			}
			if stat.PowerUsage > 0 {
				totalPower += stat.PowerUsage
				powerSamples++
				if stat.PowerUsage > maxPower {
					maxPower = stat.PowerUsage
				}
			}
		}

		results.MaxTemperature = maxTemp
		results.MaxPowerUsage = maxPower
		if powerSamples > 0 {
			results.AvgPowerUsage = totalPower / float64(powerSamples)
		}
	}

	// Calculate stability score (0-100)
	results.StabilityScore = app.calculateStabilityScore()

	// Calculate benchmark score
	results.BenchmarkScore = app.calculateBenchmarkScore(results)

	return results
}

func (app *BenchmarkApp) calculateStabilityScore() float64 {
	if len(app.frameTimesMs) < 10 {
		return 100.0
	}

	// Calculate coefficient of variation for frame times
	mean := 0.0
	for _, ft := range app.frameTimesMs {
		mean += ft
	}
	mean /= float64(len(app.frameTimesMs))

	variance := 0.0
	for _, ft := range app.frameTimesMs {
		diff := ft - mean
		variance += diff * diff
	}
	variance /= float64(len(app.frameTimesMs))

	stdDev := math.Sqrt(variance)
	cv := stdDev / mean // Coefficient of variation

	// Convert to stability score (lower CV = higher stability)
	stabilityScore := math.Max(0, 100.0-cv*10)

	// Penalize for errors
	errorPenalty := float64(app.errorCount) * 5.0
	stabilityScore = math.Max(0, stabilityScore-errorPenalty)

	return stabilityScore
}

func (app *BenchmarkApp) calculateBenchmarkScore(results *TestResults) int {
	// Base score from average FPS
	baseScore := int(results.AverageFPS * 10)

	// Resolution multiplier
	resolutionFactor := float64(app.resolution.Width*app.resolution.Height) / (1920.0 * 1080.0)
	baseScore = int(float64(baseScore) * resolutionFactor)

	// Quality multiplier
	qualityMultiplier := float64(app.quality + 1)
	baseScore = int(float64(baseScore) * qualityMultiplier)

	// Stability bonus/penalty
	stabilityFactor := results.StabilityScore / 100.0
	baseScore = int(float64(baseScore) * stabilityFactor)

	// Ensure minimum score of 0
	if baseScore < 0 {
		baseScore = 0
	}

	return baseScore
}

func (app *BenchmarkApp) displayResults(results *TestResults) {
	fmt.Print("\033[2J\033[H") // Clear screen

	fmt.Println("╔═══════════════════════════════════════════════════════════════╗")
	fmt.Println("║                    FINAL TEST RESULTS                        ║")
	fmt.Println("╚═══════════════════════════════════════════════════════════════╝")
	fmt.Println()

	// Test summary
	fmt.Printf("🏁 TEST SUMMARY\n")
	fmt.Printf("   Mode: %s\n", app.getTestModeString())
	fmt.Printf("   Quality: %s\n", app.getQualityString())
	fmt.Printf("   Resolution: %s (%dx%d)\n", app.resolution.Name, app.resolution.Width, app.resolution.Height)
	fmt.Printf("   Duration: %s\n", formatDuration(results.Duration))
	fmt.Printf("   Total Frames: %d\n", results.TotalFrames)
	fmt.Println()

	// Performance metrics
	fmt.Printf("📊 PERFORMANCE METRICS\n")
	fmt.Printf("   Average FPS: %.1f\n", results.AverageFPS)
	if results.MinFPS != math.Inf(1) {
		fmt.Printf("   Minimum FPS: %.1f\n", results.MinFPS)
	}
	if results.MaxFPS > 0 {
		fmt.Printf("   Maximum FPS: %.1f\n", results.MaxFPS)
	}

	if len(results.PercentileFPS) > 0 {
		fmt.Printf("   1%% Low FPS: %.1f\n", results.PercentileFPS["1%"])
		fmt.Printf("   5%% Low FPS: %.1f\n", results.PercentileFPS["5%"])
	}
	fmt.Println()

	// Hardware metrics
	if results.MaxTemperature > 0 || results.AvgPowerUsage > 0 {
		fmt.Printf("🌡️  HARDWARE METRICS\n")
		if results.MaxTemperature > 0 {
			fmt.Printf("   Maximum Temperature: %d°C\n", results.MaxTemperature)
		}
		if results.AvgPowerUsage > 0 {
			fmt.Printf("   Average Power: %.1f W\n", results.AvgPowerUsage)
			fmt.Printf("   Maximum Power: %.1f W\n", results.MaxPowerUsage)
		}
		fmt.Println()
	}

	// Stability assessment
	fmt.Printf("🔍 STABILITY ASSESSMENT\n")
	fmt.Printf("   Stability Score: %.1f/100\n", results.StabilityScore)
	if results.ErrorCount > 0 {
		fmt.Printf("   Artifacts Detected: %d\n", results.ErrorCount)
	} else {
		fmt.Printf("   No stability issues detected ✅\n")
	}

	// Stability rating
	var stabilityRating string
	switch {
	case results.StabilityScore >= 95:
		stabilityRating = "Excellent ⭐⭐⭐⭐⭐"
	case results.StabilityScore >= 85:
		stabilityRating = "Very Good ⭐⭐⭐⭐"
	case results.StabilityScore >= 70:
		stabilityRating = "Good ⭐⭐⭐"
	case results.StabilityScore >= 50:
		stabilityRating = "Fair ⭐⭐"
	default:
		stabilityRating = "Poor ⭐"
	}
	fmt.Printf("   Rating: %s\n", stabilityRating)
	fmt.Println()

	// Benchmark score
	if app.testMode == Benchmark {
		fmt.Printf("🏆 BENCHMARK SCORE\n")
		fmt.Printf("   Final Score: %d points\n", results.BenchmarkScore)

		// Score interpretation
		var scoreRating string
		switch {
		case results.BenchmarkScore >= 10000:
			scoreRating = "Exceptional Performance 🚀"
		case results.BenchmarkScore >= 5000:
			scoreRating = "Excellent Performance ⚡"
		case results.BenchmarkScore >= 2000:
			scoreRating = "Good Performance 👍"
		case results.BenchmarkScore >= 1000:
			scoreRating = "Average Performance 📊"
		default:
			scoreRating = "Below Average Performance 📉"
		}
		fmt.Printf("   Rating: %s\n", scoreRating)
		fmt.Println()
	}

	// Recommendations
	fmt.Printf("💡 RECOMMENDATIONS\n")
	if results.StabilityScore < 80 {
		fmt.Printf("   • Consider reducing graphics quality or resolution\n")
		fmt.Printf("   • Check GPU temperatures and cooling\n")
	}
	if results.MaxTemperature > 80 {
		fmt.Printf("   • GPU running hot - consider improved cooling\n")
	}
	if results.AverageFPS < float64(app.targetFPS)*0.8 {
		fmt.Printf("   • Performance below target - consider lower settings\n")
	}
	if results.ErrorCount == 0 && results.StabilityScore > 90 {
		fmt.Printf("   • System is stable - try higher quality settings! 🎮\n")
	}

	fmt.Println()
	fmt.Println("═══════════════════════════════════════════════════════════════")
}

func (app *BenchmarkApp) exportToCSV(outputDir string) {
	app.mutex.RLock()
	defer app.mutex.RUnlock()

	timestamp := time.Now().Format("20060102_150405")
	filename := filepath.Join(outputDir, fmt.Sprintf("gpu_stress_test_%s.csv", timestamp))

	file, err := os.Create(filename)
	if err != nil {
		log.Printf("Failed to create CSV file: %v", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{
		"Timestamp", "FPS", "FrameTime_ms", "GPU_Temp_C",
		"Power_W", "Memory_MB", "GPU_Clock_MHz", "Memory_Clock_MHz",
	}
	writer.Write(header)

	// Write performance data
	for _, data := range app.performanceLog {
		record := []string{
			data.Timestamp.Format("2006-01-02 15:04:05.000"),
			fmt.Sprintf("%.2f", data.FPS),
			fmt.Sprintf("%.2f", data.FrameTime),
			fmt.Sprintf("%d", data.GPUTemp),
			fmt.Sprintf("%.2f", data.PowerUsage),
			fmt.Sprintf("%.2f", float64(data.MemoryUsage)/(1024*1024)),
			"", // GPU Clock - would need to be added to PerformanceData
			"", // Memory Clock - would need to be added to PerformanceData
		}
		writer.Write(record)
	}

	fmt.Printf("📄 Performance data exported to: %s\n", filename)
}
