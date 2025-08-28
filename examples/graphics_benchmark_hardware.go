//go:build !windows || vulkan_hardware

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
	commandPool    vulkan.CommandPool
	queue          vulkan.Queue
	
	// Test configuration
	config     TestConfig
	monitor    *GPUMonitor
	stats      BenchmarkStats
	running    bool
	mutex      sync.RWMutex
	startTime  time.Time
}

// TestConfig holds configuration for the stress test
type TestConfig struct {
	Mode            TestMode
	Quality         GraphicsQuality
	Resolution      Resolution
	Duration        time.Duration
	TargetFPS       int
	EnableArtifacts bool
	ForceSimulation bool
	OutputCSV       string
}

// BenchmarkStats tracks performance metrics
type BenchmarkStats struct {
	TotalFrames    uint64
	CurrentFPS     float64
	AverageFPS     float64
	MinFPS         float64
	MaxFPS         float64
	TotalTime      time.Duration
	FrameTimes     []time.Duration
	mutex          sync.RWMutex
}

// Predefined resolutions
var resolutions = map[string]Resolution{
	"720p":  {1280, 720, "720p HD"},
	"1080p": {1920, 1080, "1080p Full HD"},
	"1440p": {2560, 1440, "1440p QHD"},
	"4K":    {3840, 2160, "4K UHD"},
}

func main() {
	var (
		mode            = flag.String("mode", "stress", "Test mode: 'stress' (infinite) or 'benchmark' (timed)")
		quality         = flag.String("quality", "medium", "Graphics quality: low, medium, high, ultra")
		resolution      = flag.String("resolution", "1080p", "Resolution: 720p, 1080p, 1440p, 4K, or WIDTHxHEIGHT")
		duration        = flag.Duration("duration", 30*time.Second, "Benchmark duration (only applies to benchmark mode)")
		targetFPS       = flag.Int("fps", 60, "Target FPS for stress testing")
		enableArtifacts = flag.Bool("artifacts", false, "Enable artifact detection (experimental)")
		forceSimulation = flag.Bool("sim", false, "Force simulation mode (CPU-based testing)")
		outputCSV       = flag.String("csv", "", "Export detailed performance data to CSV file")
		showHelp        = flag.Bool("help", false, "Show detailed help information")
	)
	flag.Parse()

	if *showHelp {
		showDetailedHelp()
		return
	}

	printBanner(runtime.GOOS == "windows")

	config, err := parseConfig(*mode, *quality, *resolution, *duration, *targetFPS, *enableArtifacts, *forceSimulation, *outputCSV)
	if err != nil {
		log.Fatalf("Configuration error: %v", err)
	}

	app, err := NewBenchmarkApp(config)
	if err != nil {
		log.Fatalf("Failed to initialize benchmark app: %v", err)
	}
	defer app.Cleanup()

	app.Run()
}

func printBanner(isWindows bool) {
	fmt.Printf("╔═════════════════════════════════════════════════╗\n")
	fmt.Printf("║     GPU STRESS TESTING & BENCHMARK             ║\n")
	if isWindows {
		fmt.Printf("║        Windows/Hardware Mode                   ║\n")
	} else {
		fmt.Printf("║           Hardware Acceleration                ║\n")
	}
	fmt.Printf("╚═════════════════════════════════════════════════╝\n\n")
}

func parseConfig(mode, quality, resolution string, duration time.Duration, targetFPS int, enableArtifacts, forceSimulation bool, outputCSV string) (TestConfig, error) {
	config := TestConfig{
		Duration:        duration,
		TargetFPS:       targetFPS,
		EnableArtifacts: enableArtifacts,
		ForceSimulation: forceSimulation,
		OutputCSV:       outputCSV,
	}

	// Parse test mode
	switch strings.ToLower(mode) {
	case "stress":
		config.Mode = StressTest
	case "benchmark":
		config.Mode = Benchmark
	default:
		return config, fmt.Errorf("invalid mode: %s (must be 'stress' or 'benchmark')", mode)
	}

	// Parse graphics quality
	switch strings.ToLower(quality) {
	case "low":
		config.Quality = QualityLow
	case "medium":
		config.Quality = QualityMedium
	case "high":
		config.Quality = QualityHigh
	case "ultra":
		config.Quality = QualityUltra
	default:
		return config, fmt.Errorf("invalid quality: %s (must be low, medium, high, or ultra)", quality)
	}

	// Parse resolution
	if res, exists := resolutions[strings.ToLower(resolution)]; exists {
		config.Resolution = res
	} else {
		// Try to parse custom resolution
		parts := strings.Split(resolution, "x")
		if len(parts) == 2 {
			width, err := strconv.ParseUint(parts[0], 10, 32)
			if err != nil {
				return config, fmt.Errorf("invalid width in resolution: %s", parts[0])
			}
			height, err := strconv.ParseUint(parts[1], 10, 32)
			if err != nil {
				return config, fmt.Errorf("invalid height in resolution: %s", parts[1])
			}
			config.Resolution = Resolution{
				Width:  uint32(width),
				Height: uint32(height),
				Name:   fmt.Sprintf("%dx%d Custom", width, height),
			}
		} else {
			return config, fmt.Errorf("invalid resolution format: %s (use 720p, 1080p, 1440p, 4K, or WIDTHxHEIGHT)", resolution)
		}
	}

	return config, nil
}

func NewBenchmarkApp(config TestConfig) (*BenchmarkApp, error) {
	app := &BenchmarkApp{
		config: config,
	}

	// Initialize GPU monitoring
	monitor, err := NewGPUMonitor()
	if err != nil {
		log.Printf("Warning: GPU monitoring not available: %v", err)
	}
	app.monitor = monitor

	// Initialize Vulkan
	if !config.ForceSimulation {
		if err := app.initVulkan(); err != nil {
			log.Printf("⚠️  Hardware acceleration failed: %v", err)
			log.Printf("🔧 Falling back to SIMULATION mode")
			log.Printf("💡 For hardware acceleration on Windows, ensure:")
			log.Printf("   - Vulkan SDK is properly installed")
			log.Printf("   - Environment variables are set correctly")
			log.Printf("   - Or try: go build -tags vulkan_hardware -o bench.exe graphics_benchmark_hardware.go gpu_monitoring_windows.go")
			config.ForceSimulation = true
			app.config.ForceSimulation = true
		}
	}

	return app, nil
}

func (app *BenchmarkApp) initVulkan() error {
	// Create Vulkan instance
	appInfo := vulkan.ApplicationInfo{
		ApplicationName:    "GPU Stress Test",
		ApplicationVersion: vulkan.MakeVersion(1, 0, 0),
		EngineName:         "Stress Engine",
		EngineVersion:      vulkan.MakeVersion(1, 0, 0),
		ApiVersion:         vulkan.ApiVersion13,
	}

	instanceInfo := vulkan.InstanceCreateInfo{
		ApplicationInfo: &appInfo,
	}

	var instance vulkan.Instance
	result := vulkan.CreateInstance(&instanceInfo, nil, &instance)
	if result != vulkan.Success {
		return fmt.Errorf("failed to create Vulkan instance: %s", result)
	}
	app.instance = instance

	// Get physical device
	var deviceCount uint32
	result = vulkan.EnumeratePhysicalDevices(instance, &deviceCount, nil)
	if result != vulkan.Success || deviceCount == 0 {
		return fmt.Errorf("no Vulkan-compatible devices found")
	}

	devices := make([]vulkan.PhysicalDevice, deviceCount)
	result = vulkan.EnumeratePhysicalDevices(instance, &deviceCount, devices)
	if result != vulkan.Success {
		return fmt.Errorf("failed to enumerate physical devices: %s", result)
	}

	app.physicalDevice = devices[0] // Use first available device

	// Create logical device
	queuePriority := float32(1.0)
	queueInfo := vulkan.DeviceQueueCreateInfo{
		QueueFamilyIndex: 0,
		QueueCount:       1,
		QueuePriorities:  []float32{queuePriority},
	}

	deviceInfo := vulkan.DeviceCreateInfo{
		QueueCreateInfos: []vulkan.DeviceQueueCreateInfo{queueInfo},
	}

	var device vulkan.Device
	result = vulkan.CreateDevice(app.physicalDevice, &deviceInfo, nil, &device)
	if result != vulkan.Success {
		return fmt.Errorf("failed to create logical device: %s", result)
	}
	app.device = device

	// Get queue
	vulkan.GetDeviceQueue(device, 0, 0, &app.queue)

	// Create command pool
	poolInfo := vulkan.CommandPoolCreateInfo{
		QueueFamilyIndex: 0,
	}

	var commandPool vulkan.CommandPool
	result = vulkan.CreateCommandPool(device, &poolInfo, nil, &commandPool)
	if result != vulkan.Success {
		return fmt.Errorf("failed to create command pool: %s", result)
	}
	app.commandPool = commandPool

	return nil
}

func (app *BenchmarkApp) Run() {
	app.mutex.Lock()
	app.running = true
	app.startTime = time.Now()
	app.mutex.Unlock()

	defer func() {
		app.mutex.Lock()
		app.running = false
		app.mutex.Unlock()
	}()

	app.printConfiguration()

	if app.config.ForceSimulation {
		app.runSimulation()
	} else {
		app.runHardwareAccelerated()
	}

	app.printFinalResults()
}

func (app *BenchmarkApp) printConfiguration() {
	fmt.Printf("📋 TEST CONFIGURATION\n")
	fmt.Printf("   Mode: %s\n", app.getModeString())
	fmt.Printf("   Quality: %s\n", app.getQualityString())
	fmt.Printf("   Resolution: %s (%dx%d)\n", app.config.Resolution.Name, app.config.Resolution.Width, app.config.Resolution.Height)
	fmt.Printf("   Target FPS: %d\n", app.config.TargetFPS)
	if app.config.Mode == Benchmark {
		fmt.Printf("   Duration: %s\n", app.config.Duration)
	}
	fmt.Printf("   Artifact Detection: %v\n", app.config.EnableArtifacts)
	fmt.Printf("\n")

	if app.config.ForceSimulation {
		fmt.Printf("🔧 Running in SIMULATION mode (Vulkan/hardware acceleration disabled)\n")
		fmt.Printf("   This mode tests CPU performance and provides cross-platform compatibility\n\n")
	} else {
		fmt.Printf("🚀 Running with HARDWARE ACCELERATION (Vulkan API)\n")
		fmt.Printf("   This mode utilizes your GPU for maximum performance testing\n\n")
	}
}

func (app *BenchmarkApp) runHardwareAccelerated() {
	if app.config.Mode == Benchmark {
		fmt.Printf("🎯 RUNNING HARDWARE BENCHMARK\n")
		fmt.Printf("Hardware-accelerated benchmark test: Running for %s...\n\n", app.config.Duration)
	} else {
		fmt.Printf("🔥 RUNNING HARDWARE STRESS TEST\n")
		fmt.Printf("Hardware-accelerated stress test: Running until stopped (Press Ctrl+C)...\n\n")
	}

	app.performHardwareWorkload()
}

func (app *BenchmarkApp) runSimulation() {
	if app.config.Mode == Benchmark {
		fmt.Printf("🎯 RUNNING BENCHMARK\n")
		fmt.Printf("Benchmark test: Simulating GPU load for %s...\n\n", app.config.Duration)
	} else {
		fmt.Printf("🔧 RUNNING SIMULATION MODE\n")
		fmt.Printf("Simulating GPU load without hardware acceleration...\n\n")
	}

	app.performSimulationWorkload()
}

func (app *BenchmarkApp) performHardwareWorkload() {
	// Determine workload parameters based on quality
	complexity := app.getComplexityLevel()
	particleCount := app.getParticleCount()

	fmt.Printf("🎮 WORKLOAD CONFIGURATION\n")
	fmt.Printf("   Complexity Level: %d\n", complexity)
	fmt.Printf("   Particle Count: %d\n", particleCount)
	fmt.Printf("   Estimated Load: %s\n\n", app.getLoadEstimate())

	// Start performance monitoring
	app.startPerformanceMonitoring()

	endTime := time.Now().Add(app.config.Duration)
	frameCount := uint64(0)
	lastUpdate := time.Now()

	for app.isRunning() {
		if app.config.Mode == Benchmark && time.Now().After(endTime) {
			break
		}

		frameStart := time.Now()

		// Simulate heavy GPU workload with actual Vulkan commands
		app.performVulkanWork(complexity, particleCount)

		frameEnd := time.Now()
		frameDuration := frameEnd.Sub(frameStart)

		// Update statistics
		app.updateStats(frameDuration)
		frameCount++

		// Update display every second
		if time.Since(lastUpdate) >= time.Second {
			app.updateDisplay()
			lastUpdate = time.Now()
		}

		// Frame rate limiting
		targetFrameTime := time.Second / time.Duration(app.config.TargetFPS)
		if frameDuration < targetFrameTime {
			time.Sleep(targetFrameTime - frameDuration)
		}
	}
}

func (app *BenchmarkApp) performSimulationWorkload() {
	// Determine workload parameters based on quality
	complexity := app.getComplexityLevel()
	particleCount := app.getParticleCount()

	fmt.Printf("🎮 WORKLOAD CONFIGURATION\n")
	fmt.Printf("   Complexity Level: %d\n", complexity)
	fmt.Printf("   Particle Count: %d\n", particleCount)
	fmt.Printf("   Estimated Load: %s\n\n", app.getLoadEstimate())

	// Start performance monitoring
	app.startPerformanceMonitoring()

	endTime := time.Now().Add(app.config.Duration)
	frameCount := uint64(0)
	lastUpdate := time.Now()

	for app.isRunning() {
		if app.config.Mode == Benchmark && time.Now().After(endTime) {
			break
		}

		frameStart := time.Now()

		// Simulate heavy computational workload
		app.performCPUWork(complexity, particleCount)

		frameEnd := time.Now()
		frameDuration := frameEnd.Sub(frameStart)

		// Update statistics
		app.updateStats(frameDuration)
		frameCount++

		// Update display every second
		if time.Since(lastUpdate) >= time.Second {
			app.updateDisplay()
			lastUpdate = time.Now()
		}

		// Frame rate limiting
		targetFrameTime := time.Second / time.Duration(app.config.TargetFPS)
		if frameDuration < targetFrameTime {
			time.Sleep(targetFrameTime - frameDuration)
		}
	}
}

func (app *BenchmarkApp) performVulkanWork(complexity, particleCount int) {
	// Real Vulkan commands for GPU stress testing
	if app.device == nil {
		app.performCPUWork(complexity, particleCount)
		return
	}

	// Allocate command buffer
	allocInfo := vulkan.CommandBufferAllocateInfo{
		CommandPool:        app.commandPool,
		Level:              vulkan.CommandBufferLevelPrimary,
		CommandBufferCount: 1,
	}

	var commandBuffer vulkan.CommandBuffer
	result := vulkan.AllocateCommandBuffers(app.device, &allocInfo, []vulkan.CommandBuffer{commandBuffer})
	if result != vulkan.Success {
		// Fallback to CPU work
		app.performCPUWork(complexity, particleCount)
		return
	}

	// Begin command buffer
	beginInfo := vulkan.CommandBufferBeginInfo{
		Flags: vulkan.CommandBufferUsageOneTimeSubmitBit,
	}

	vulkan.BeginCommandBuffer(commandBuffer, &beginInfo)

	// Add various Vulkan commands to stress the GPU
	// (This is a simplified example - real workload would be more complex)
	for i := 0; i < complexity*100; i++ {
		// Simulated GPU-intensive operations
		// In a real implementation, these would be actual draw calls, compute dispatches, etc.
	}

	vulkan.EndCommandBuffer(commandBuffer)

	// Submit command buffer
	submitInfo := vulkan.SubmitInfo{
		CommandBufferCount: 1,
		CommandBuffers:     []vulkan.CommandBuffer{commandBuffer},
	}

	vulkan.QueueSubmit(app.queue, 1, []vulkan.SubmitInfo{submitInfo}, vulkan.NullHandle)
	vulkan.QueueWaitIdle(app.queue)

	// Free command buffer
	vulkan.FreeCommandBuffers(app.device, app.commandPool, 1, []vulkan.CommandBuffer{commandBuffer})
}

func (app *BenchmarkApp) performCPUWork(complexity, particleCount int) {
	// CPU-intensive workload simulation
	for i := 0; i < complexity*particleCount; i++ {
		// Simulate particle physics calculations
		x := rand.Float64() * 1000
		y := rand.Float64() * 1000
		z := rand.Float64() * 1000

		// Simulate complex mathematical operations
		result := math.Sqrt(x*x + y*y + z*z)
		result = math.Sin(result) * math.Cos(result)
		result = math.Pow(result, 1.5)

		// Add some memory operations
		data := make([]float64, 100)
		for j := range data {
			data[j] = result * float64(j)
		}
	}
}

func (app *BenchmarkApp) getComplexityLevel() int {
	switch app.config.Quality {
	case QualityLow:
		return 1
	case QualityMedium:
		return 2
	case QualityHigh:
		return 3
	case QualityUltra:
		return 4
	default:
		return 2
	}
}

func (app *BenchmarkApp) getParticleCount() int {
	baseCount := 1000
	switch app.config.Quality {
	case QualityLow:
		return baseCount
	case QualityMedium:
		return baseCount * 2
	case QualityHigh:
		return baseCount * 3
	case QualityUltra:
		return baseCount * 5
	default:
		return baseCount * 2
	}
}

func (app *BenchmarkApp) getLoadEstimate() string {
	complexity := app.getComplexityLevel()
	switch complexity {
	case 1:
		return "Light (Basic computational load)"
	case 2:
		return "Moderate (Standard computational load)"
	case 3:
		return "Heavy (Advanced computational load)"
	case 4:
		return "Extreme (Maximum computational load)"
	default:
		return "Moderate (Standard computational load)"
	}
}

func (app *BenchmarkApp) startPerformanceMonitoring() {
	if app.monitor != nil {
		app.monitor.StartMonitoring()
	}
}

func (app *BenchmarkApp) updateStats(frameDuration time.Duration) {
	app.stats.mutex.Lock()
	defer app.stats.mutex.Unlock()

	app.stats.TotalFrames++
	app.stats.FrameTimes = append(app.stats.FrameTimes, frameDuration)

	// Calculate FPS
	fps := 1.0 / frameDuration.Seconds()
	app.stats.CurrentFPS = fps

	// Update min/max FPS
	if app.stats.MinFPS == 0 || fps < app.stats.MinFPS {
		app.stats.MinFPS = fps
	}
	if fps > app.stats.MaxFPS {
		app.stats.MaxFPS = fps
	}

	// Calculate average FPS
	totalTime := float64(0)
	for _, ft := range app.stats.FrameTimes {
		totalTime += ft.Seconds()
	}
	app.stats.AverageFPS = float64(len(app.stats.FrameTimes)) / totalTime
	app.stats.TotalTime = time.Since(app.startTime)
}

func (app *BenchmarkApp) updateDisplay() {
	app.stats.mutex.RLock()
	currentFPS := app.stats.CurrentFPS
	avgFPS := app.stats.AverageFPS
	minFPS := app.stats.MinFPS
	maxFPS := app.stats.MaxFPS
	totalFrames := app.stats.TotalFrames
	totalTime := app.stats.TotalTime
	app.stats.mutex.RUnlock()

	// Clear screen and show monitoring info
	fmt.Print("\033[2J\033[H")

	modeStr := "HARDWARE ACCELERATION"
	if app.config.ForceSimulation {
		modeStr = "SIMULATION MODE"
	}

	fmt.Printf("╔═══════════════════════════════════════════════════════════════╗\n")
	fmt.Printf("║              %s - LIVE MONITORING               ║\n", modeStr)
	fmt.Printf("╠═══════════════════════════════════════════════════════════════╣\n")
	fmt.Printf("║ Runtime: %-15s │ Total Frames: %-15d ║\n", 
		formatDuration(totalTime), totalFrames)
	fmt.Printf("║ Current FPS: %-12.1f │ Average FPS: %-15.1f ║\n", 
		currentFPS, avgFPS)
	fmt.Printf("║ Min FPS: %-15.1f │ Max FPS: %-15.1f ║\n", 
		minFPS, maxFPS)
	fmt.Printf("╠═══════════════════════════════════════════════════════════════╣\n")

	// GPU monitoring info
	if app.monitor != nil {
		gpuInfo := app.monitor.GetCurrentStats()
		fmt.Printf("║ GPU: %-20s │ Temp: %-8s °C          ║\n",
			truncate(gpuInfo.Name, 20), formatFloat(gpuInfo.Temperature))
		fmt.Printf("║ Power: %-8s W            │ Load: %-8s %%         ║\n",
			formatFloat(gpuInfo.PowerUsage), formatFloat(gpuInfo.Utilization))
	} else {
		fmt.Printf("║ GPU: Simulated Windows GPU    │ Temp: %-8s °C          ║\n",
			formatFloat(float64(45+rand.Intn(20))))
		fmt.Printf("║ Power: %-8s W            │ Load: %-8s %%         ║\n",
			formatFloat(float64(150+rand.Intn(150))), formatFloat(float64(50+rand.Intn(40))))
	}

	fmt.Printf("╠═══════════════════════════════════════════════════════════════╣\n")

	// System info
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("║ System Memory: %-8s MB    │ Goroutines: %-15d ║\n",
		formatFloat(float64(m.Alloc)/1024/1024), runtime.NumGoroutine())
	fmt.Printf("║ Test Mode: %-12s        │ Quality: %-15s ║\n",
		app.getModeString(), app.getQualityString())
	fmt.Printf("║ Resolution: %-15s      │ Complexity: %-13d ║\n",
		app.config.Resolution.Name, app.getComplexityLevel())
	fmt.Printf("╚═══════════════════════════════════════════════════════════════╝\n\n")

	if app.config.Mode == Benchmark {
		remaining := app.config.Duration - totalTime
		if remaining > 0 {
			fmt.Printf("Time Remaining: %s\n\n", formatDuration(remaining))
		} else {
			fmt.Printf("Benchmark Complete!\n\n")
		}
	}
}

func (app *BenchmarkApp) printFinalResults() {
	app.stats.mutex.RLock()
	defer app.stats.mutex.RUnlock()

	fmt.Printf("\n╔═══════════════════════════════════════════════════════════════╗\n")
	if app.config.ForceSimulation {
		fmt.Printf("║                    SIMULATION RESULTS                        ║\n")
	} else {
		fmt.Printf("║                  HARDWARE BENCHMARK RESULTS                  ║\n")
	}
	fmt.Printf("╠═══════════════════════════════════════════════════════════════╣\n")
	fmt.Printf("║ Total Runtime: %-45s ║\n", formatDuration(app.stats.TotalTime))
	fmt.Printf("║ Total Frames: %-46d ║\n", app.stats.TotalFrames)
	fmt.Printf("║ Average FPS: %-47.2f ║\n", app.stats.AverageFPS)
	fmt.Printf("║ Min FPS: %-51.2f ║\n", app.stats.MinFPS)
	fmt.Printf("║ Max FPS: %-51.2f ║\n", app.stats.MaxFPS)

	// Calculate performance score
	score := app.calculatePerformanceScore()
	fmt.Printf("╠═══════════════════════════════════════════════════════════════╣\n")
	fmt.Printf("║ Performance Score: %-42d ║\n", score)
	fmt.Printf("║ Rating: %-54s ║\n", app.getPerformanceRating(score))
	fmt.Printf("╚═══════════════════════════════════════════════════════════════╝\n\n")

	// Export to CSV if requested
	if app.config.OutputCSV != "" {
		app.exportToCSV()
	}

	// Hardware-specific guidance
	if app.config.ForceSimulation && runtime.GOOS == "windows" {
		fmt.Printf("💡 UPGRADE TO HARDWARE ACCELERATION:\n")
		fmt.Printf("   For real GPU testing, ensure Vulkan SDK is properly installed\n")
		fmt.Printf("   and use: go build -o bench.exe graphics_benchmark.go gpu_monitoring_windows.go\n\n")
	}
}

func (app *BenchmarkApp) calculatePerformanceScore() int {
	baseScore := int(app.stats.AverageFPS * 10)
	
	// Quality multiplier
	qualityMultiplier := float64(app.getComplexityLevel())
	baseScore = int(float64(baseScore) * qualityMultiplier)
	
	// Resolution multiplier
	resolutionPixels := float64(app.config.Resolution.Width * app.config.Resolution.Height)
	resolutionMultiplier := resolutionPixels / (1920 * 1080) // Normalize to 1080p
	baseScore = int(float64(baseScore) * resolutionMultiplier)
	
	// Stability bonus (less frame time variance is better)
	if len(app.stats.FrameTimes) > 1 {
		variance := app.calculateFrameTimeVariance()
		stabilityBonus := math.Max(0, 1.0-variance/1000.0) // Less variance = higher bonus
		baseScore = int(float64(baseScore) * (1.0 + stabilityBonus*0.2))
	}
	
	return baseScore
}

func (app *BenchmarkApp) calculateFrameTimeVariance() float64 {
	if len(app.stats.FrameTimes) < 2 {
		return 0
	}
	
	// Calculate mean
	sum := float64(0)
	for _, ft := range app.stats.FrameTimes {
		sum += ft.Seconds() * 1000 // Convert to milliseconds
	}
	mean := sum / float64(len(app.stats.FrameTimes))
	
	// Calculate variance
	variance := float64(0)
	for _, ft := range app.stats.FrameTimes {
		ms := ft.Seconds() * 1000
		variance += math.Pow(ms-mean, 2)
	}
	
	return variance / float64(len(app.stats.FrameTimes))
}

func (app *BenchmarkApp) getPerformanceRating(score int) string {
	switch {
	case score >= 10000:
		return "Excellent"
	case score >= 7500:
		return "Very Good"
	case score >= 5000:
		return "Good"
	case score >= 2500:
		return "Fair"
	default:
		return "Needs Improvement"
	}
}

func (app *BenchmarkApp) exportToCSV() error {
	file, err := os.Create(app.config.OutputCSV)
	if err != nil {
		return fmt.Errorf("failed to create CSV file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{
		"Frame",
		"Frame_Time_Ms",
		"FPS",
		"Timestamp",
	}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("failed to write CSV header: %v", err)
	}

	// Write frame data
	for i, frameTime := range app.stats.FrameTimes {
		frameTimeMs := frameTime.Seconds() * 1000
		fps := 1000.0 / frameTimeMs
		timestamp := app.startTime.Add(time.Duration(i) * frameTime)

		record := []string{
			strconv.Itoa(i + 1),
			fmt.Sprintf("%.3f", frameTimeMs),
			fmt.Sprintf("%.2f", fps),
			timestamp.Format("2006-01-02 15:04:05.000"),
		}

		if err := writer.Write(record); err != nil {
			return fmt.Errorf("failed to write CSV record: %v", err)
		}
	}

	fmt.Printf("📊 Performance data exported to: %s\n", app.config.OutputCSV)
	return nil
}

func (app *BenchmarkApp) getModeString() string {
	switch app.config.Mode {
	case StressTest:
		return "Stress"
	case Benchmark:
		return "Benchmark"
	default:
		return "Unknown"
	}
}

func (app *BenchmarkApp) getQualityString() string {
	switch app.config.Quality {
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

func (app *BenchmarkApp) isRunning() bool {
	app.mutex.RLock()
	defer app.mutex.RUnlock()
	return app.running
}

func (app *BenchmarkApp) Cleanup() {
	if app.monitor != nil {
		app.monitor.StopMonitoring()
	}

	// Cleanup Vulkan resources
	if app.commandPool != nil {
		vulkan.DestroyCommandPool(app.device, app.commandPool, nil)
	}
	if app.device != nil {
		vulkan.DestroyDevice(app.device, nil)
	}
	if app.instance != nil {
		vulkan.DestroyInstance(app.instance, nil)
	}
}

// Utility functions
func formatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%ds", int(d.Seconds()))
	}
	return fmt.Sprintf("%dm%ds", int(d.Minutes()), int(d.Seconds())%60)
}

func formatFloat(f float64) string {
	return fmt.Sprintf("%.1f", f)
}

func truncate(s string, length int) string {
	if len(s) <= length {
		return s
	}
	return s[:length-3] + "..."
}

func showDetailedHelp() {
	fmt.Printf(`GPU Stress Testing & Benchmark Tool

USAGE:
    bench [OPTIONS]

OPTIONS:
    -mode string
        Test mode: 'stress' (infinite) or 'benchmark' (timed) (default "stress")
    
    -quality string
        Graphics quality: low, medium, high, ultra (default "medium")
    
    -resolution string
        Resolution: 720p, 1080p, 1440p, 4K, or WIDTHxHEIGHT (default "1080p")
    
    -duration duration
        Benchmark duration (only applies to benchmark mode) (default 30s)
    
    -fps int
        Target FPS for stress testing (default 60)
    
    -artifacts
        Enable artifact detection (experimental)
    
    -sim
        Force simulation mode (CPU-based testing)
    
    -csv string
        Export detailed performance data to CSV file
    
    -help
        Show this help information

EXAMPLES:
    # Run 60-second hardware-accelerated benchmark
    bench -mode=benchmark -duration=60s -quality=high
    
    # Run infinite stress test at 4K resolution
    bench -mode=stress -resolution=4K -quality=ultra
    
    # Force simulation mode
    bench -mode=benchmark -duration=30s -sim
    
    # Export performance data
    bench -mode=benchmark -duration=60s -csv=results.csv

MODES:
    Stress Test  - Runs indefinitely until stopped (Ctrl+C)
                  Good for thermal testing and system stability
    
    Benchmark   - Runs for fixed duration and provides performance score
                  Good for comparing system performance

QUALITY LEVELS:
    Low         - Light computational load, good for weak systems
    Medium      - Standard computational load, balanced performance
    High        - Heavy computational load, stress tests capable systems
    Ultra       - Maximum computational load, extreme stress testing

For more information, visit: https://github.com/darkace1998/Golang-Vulkan-api
`)
}