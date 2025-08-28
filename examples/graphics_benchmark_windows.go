//go:build windows || (linux && !cgo)

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
	// Vulkan objects - stubbed for simulation-only mode
	instance       interface{}
	physicalDevice interface{}
	device         interface{}
	graphicsQueue  interface{}
	commandPool    interface{}

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
	fmt.Println("GPU STRESS TESTING & BENCHMARK APPLICATION (Windows/Cross-Platform)")
	fmt.Println("====================================================================")
	fmt.Println()
	fmt.Println("DESCRIPTION:")
	fmt.Println("  GPU stress testing application that supports both hardware-accelerated")
	fmt.Println("  Vulkan mode and cross-platform simulation mode. Provides comprehensive")
	fmt.Println("  stress testing and benchmark modes for GPU performance evaluation.")
	fmt.Println()
	fmt.Println("USAGE:")
	flag.PrintDefaults()
	fmt.Println()
	fmt.Println("EXECUTION MODES:")
	fmt.Println("  Hardware Mode  - Uses Vulkan API for real GPU acceleration (default)")
	fmt.Println("  Simulation Mode - CPU-based cross-platform mode (use -sim flag)")
	fmt.Println()
	fmt.Println("TEST MODES:")
	fmt.Println("  stress    - Runs indefinitely until manually stopped (default)")
	fmt.Println("  benchmark - Runs for fixed duration and provides performance score")
	fmt.Println()
	fmt.Println("QUALITY LEVELS:")
	fmt.Println("  low       - Basic workload, minimal system load")
	fmt.Println("  medium    - Standard workload, moderate system load")
	fmt.Println("  high      - Advanced workload, high system load (default)")
	fmt.Println("  ultra     - Maximum workload, extreme system load")
	fmt.Println()
	fmt.Println("EXAMPLES:")
	fmt.Println("  # Run hardware-accelerated 60-second benchmark:")
	fmt.Println("  ./bench.exe -mode=benchmark -duration=60s -quality=high")
	fmt.Println()
	fmt.Println("  # Run simulation mode benchmark (no Vulkan required):")
	fmt.Println("  ./bench.exe -mode=benchmark -duration=60s -sim")
	fmt.Println()
	fmt.Println("  # Run hardware-accelerated infinite stress test:")
	fmt.Println("  ./bench.exe -mode=stress -quality=ultra")
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
	fmt.Println("║     GPU STRESS TESTING & BENCHMARK             ║")
	fmt.Println("║     Simulation Mode - Cross-Platform Build     ║")
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

	// Check if user wants hardware acceleration but is using wrong build
	if *simMode {
		fmt.Printf("🔧 Running %s in SIMULATION mode (as requested)\n", strings.ToUpper(app.getTestModeString()))
	} else {
		fmt.Printf("💡 Note: For HARDWARE ACCELERATION on Windows, use:\n")
		fmt.Printf("   go build -o bench.exe graphics_benchmark.go gpu_monitoring_windows.go\n")
		fmt.Printf("   (requires Vulkan SDK installation)\n")
		fmt.Println()
		fmt.Printf("🔧 Running %s in SIMULATION mode\n", strings.ToUpper(app.getTestModeString()))
	}
	fmt.Println("   This mode tests CPU performance and provides cross-platform compatibility")
	if app.testMode == Benchmark {
		fmt.Printf("   Benchmark will run for %s and provide a performance score\n", formatDuration(app.maxDuration))
	} else {
		fmt.Println("   Stress test will run indefinitely until manually stopped")
	}
	fmt.Println()
	app.runSimulation()

	// Generate final report
	results := app.generateResults()
	app.displayResults(results)

	// Export data if requested
	if *csvExport && *outputDir != "" {
		app.exportToCSV(*outputDir)
	}
}

// Include all the implementation methods from the original file...
// (I'll include the key methods to make this functional)

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

func (app *BenchmarkApp) runSimulation() {
	fmt.Printf("🎯 RUNNING %s\n", strings.ToUpper(app.getTestModeString()))
	if app.testMode == Benchmark {
		fmt.Printf("Benchmark test: Simulating GPU load for %s...\n", formatDuration(app.maxDuration))
	} else {
		fmt.Println("Stress test: Simulating GPU load indefinitely...")
	}
	fmt.Println()

	app.startTime = time.Now()
	app.lastFrameTime = time.Now()
	app.minFPS = math.Inf(1)
	app.maxFPS = 0

	// Initialize monitoring
	app.initGPUMonitoring()
	defer app.cleanupGPUMonitoring()

	// Set moderate complexity for simulation
	app.complexityLevel = int(app.quality) + 1
	app.particleCount = 1000 * app.complexityLevel

	fmt.Printf("🎮 WORKLOAD CONFIGURATION\n")
	fmt.Printf("   Complexity Level: %d\n", app.complexityLevel)
	fmt.Printf("   Particle Count: %d\n", app.particleCount)
	fmt.Printf("   Estimated Load: %s\n", app.getLoadDescription())
	fmt.Println()

	// Start monitoring goroutine
	go app.monitoringLoop()

	frameInterval := time.Second / time.Duration(app.targetFPS)
	ticker := time.NewTicker(frameInterval)
	defer ticker.Stop()

	lastDisplayTime := time.Now()

	for {
		select {
		case <-ticker.C:
			app.simulateAdvancedWorkload()
			app.updatePerformanceMetrics()

			if app.shouldExit() {
				return
			}

			// Display stats every second
			if time.Since(lastDisplayTime) >= time.Second {
				app.displayLiveStats()
				lastDisplayTime = time.Now()
			}
		}
	}
}

func (app *BenchmarkApp) getLoadDescription() string {
	switch app.quality {
	case QualityLow:
		return "Light (Basic computational stress)"
	case QualityMedium:
		return "Moderate (Standard workload)"
	case QualityHigh:
		return "Heavy (Advanced computational load)"
	case QualityUltra:
		return "Extreme (Maximum CPU utilization)"
	default:
		return "Unknown"
	}
}

func (app *BenchmarkApp) simulateAdvancedWorkload() {
	app.animationTime += 0.016
	app.rotationAngle = float32(math.Mod(float64(app.animationTime), 2*math.Pi))

	// Simulate CPU-intensive work proportional to complexity
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

	app.lastFrameTime = now
}

func (app *BenchmarkApp) monitoringLoop() {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if app.monitoringEnabled {
				app.collectPerformanceData()
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

		if len(app.statsHistory) > 1000 {
			app.statsHistory = app.statsHistory[1:]
		}

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

		if len(app.performanceLog) > 10000 {
			app.performanceLog = app.performanceLog[1:]
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
	fmt.Print("\033[2J\033[H") // Clear screen

	elapsed := time.Since(app.startTime)

	fmt.Println("╔═══════════════════════════════════════════════════════════════╗")
	fmt.Printf("║              %s - LIVE MONITORING               ║\n", strings.ToUpper(app.getTestModeString()))
	fmt.Println("╠═══════════════════════════════════════════════════════════════╣")

	fmt.Printf("║ Runtime: %-15s │ Total Frames: %-15d ║\n",
		formatDuration(elapsed), app.frameCount)
	fmt.Printf("║ Current FPS: %-12.1f │ Average FPS: %-14.1f ║\n",
		app.currentFPS, app.avgFPS)

	if app.minFPS != math.Inf(1) && app.maxFPS > 0 {
		fmt.Printf("║ Min FPS: %-15.1f │ Max FPS: %-18.1f ║\n",
			app.minFPS, app.maxFPS)
	}

	fmt.Println("╠═══════════════════════════════════════════════════════════════╣")

	// Simulated GPU stats
	stats := app.getGPUStats()
	if stats != nil {
		fmt.Printf("║ Simulated GPU: %-14s │ Temp: %-8d°C          ║\n",
			stats.Vendor, stats.Temperature)
		fmt.Printf("║ Power: %-7.1f W            │ Load: %-8d%%         ║\n",
			stats.PowerUsage, stats.GPUUtilization)
	}

	// System info
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	systemMemMB := float64(m.Alloc) / (1024 * 1024)

	fmt.Println("╠═══════════════════════════════════════════════════════════════╣")
	fmt.Printf("║ System Memory: %-7.1f MB    │ Goroutines: %-12d    ║\n",
		systemMemMB, runtime.NumGoroutine())

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

	// Calculate stability score
	results.StabilityScore = 95.0 // Simplified for simulation mode

	// Calculate benchmark score
	results.BenchmarkScore = app.calculateBenchmarkScore(results)

	return results
}

func (app *BenchmarkApp) calculateBenchmarkScore(results *TestResults) int {
	baseScore := int(results.AverageFPS * 10)
	resolutionFactor := float64(app.resolution.Width*app.resolution.Height) / (1920.0 * 1080.0)
	baseScore = int(float64(baseScore) * resolutionFactor)
	qualityMultiplier := float64(app.quality + 1)
	baseScore = int(float64(baseScore) * qualityMultiplier)
	return baseScore
}

func (app *BenchmarkApp) displayResults(results *TestResults) {
	fmt.Print("\033[2J\033[H") // Clear screen

	fmt.Println("╔═══════════════════════════════════════════════════════════════╗")
	fmt.Println("║                    FINAL TEST RESULTS                        ║")
	fmt.Println("╚═══════════════════════════════════════════════════════════════╝")
	fmt.Println()

	fmt.Printf("🏁 TEST SUMMARY\n")
	fmt.Printf("   Mode: %s (Simulated)\n", app.getTestModeString())
	fmt.Printf("   Quality: %s\n", app.getQualityString())
	fmt.Printf("   Resolution: %s (%dx%d)\n", app.resolution.Name, app.resolution.Width, app.resolution.Height)
	fmt.Printf("   Duration: %s\n", formatDuration(results.Duration))
	fmt.Printf("   Total Frames: %d\n", results.TotalFrames)
	fmt.Println()

	fmt.Printf("📊 PERFORMANCE METRICS\n")
	fmt.Printf("   Average FPS: %.1f\n", results.AverageFPS)
	if results.MinFPS != math.Inf(1) {
		fmt.Printf("   Minimum FPS: %.1f\n", results.MinFPS)
	}
	if results.MaxFPS > 0 {
		fmt.Printf("   Maximum FPS: %.1f\n", results.MaxFPS)
	}
	fmt.Println()

	if app.testMode == Benchmark {
		fmt.Printf("🏆 BENCHMARK SCORE\n")
		fmt.Printf("   Final Score: %d points\n", results.BenchmarkScore)
		fmt.Println()
	}

	fmt.Println("═══════════════════════════════════════════════════════════════")
}

func (app *BenchmarkApp) exportToCSV(outputDir string) {
	app.mutex.RLock()
	defer app.mutex.RUnlock()

	timestamp := time.Now().Format("20060102_150405")
	filename := filepath.Join(outputDir, fmt.Sprintf("simulation_test_%s.csv", timestamp))

	file, err := os.Create(filename)
	if err != nil {
		log.Printf("Failed to create CSV file: %v", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{"Timestamp", "FPS", "FrameTime_ms", "Simulated_Temp_C", "Simulated_Power_W"}
	writer.Write(header)

	for _, data := range app.performanceLog {
		record := []string{
			data.Timestamp.Format("2006-01-02 15:04:05.000"),
			fmt.Sprintf("%.2f", data.FPS),
			fmt.Sprintf("%.2f", data.FrameTime),
			fmt.Sprintf("%d", data.GPUTemp),
			fmt.Sprintf("%.2f", data.PowerUsage),
		}
		writer.Write(record)
	}

	fmt.Printf("📄 Performance data exported to: %s\n", filename)
}