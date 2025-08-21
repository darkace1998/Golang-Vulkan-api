package main

import (
	"strings"
	"testing"
	"time"
)

func TestBenchmarkApp_InitVulkan(t *testing.T) {
	app := &BenchmarkApp{
		startTime:     time.Now(),
		lastFrameTime: time.Now(),
	}
	
	// Test Vulkan initialization (should fail gracefully without drivers)
	err := app.initVulkan()
	if err != nil {
		// This is expected in CI environments without GPU drivers
		if !strings.Contains(err.Error(), "VK_ERROR_INCOMPATIBLE_DRIVER") &&
			!strings.Contains(err.Error(), "failed to create instance") {
			t.Errorf("Unexpected error type: %v", err)
		}
	} else {
		// If Vulkan init succeeds, clean up
		app.cleanup()
	}
}

func TestBenchmarkApp_GPUMonitoring(t *testing.T) {
	app := &BenchmarkApp{}
	
	// Test GPU monitoring initialization
	app.initGPUMonitoring()
	defer app.cleanupGPUMonitoring()
	
	// Try to get GPU stats (may return nil if no NVIDIA GPU)
	stats := app.getGPUStats()
	
	// This should not panic, even if no GPU is available
	if stats != nil {
		t.Logf("GPU stats available: temp=%d°C, mem_clock=%dMHz, gpu_clock=%dMHz",
			stats.Temperature, stats.MemoryClock, stats.GraphicsClock)
	} else {
		t.Log("GPU stats not available (expected in CI environment)")
	}
}

func TestBenchmarkApp_RenderFrame(t *testing.T) {
	app := &BenchmarkApp{
		startTime:     time.Now(),
		lastFrameTime: time.Now(),
	}
	
	initialFrame := app.frameCount
	initialAngle := app.rotationAngle
	
	// Render a few frames
	for i := 0; i < 10; i++ {
		app.renderFrame()
	}
	
	// Check that frame count increased
	if app.frameCount != initialFrame+10 {
		t.Errorf("Expected frame count %d, got %d", initialFrame+10, app.frameCount)
	}
	
	// Check that rotation angle changed
	if app.rotationAngle == initialAngle {
		t.Error("Rotation angle should have changed")
	}
	
	// Check that FPS was calculated
	if app.currentFPS <= 0 {
		t.Error("Current FPS should be positive")
	}
}

func TestBenchmarkApp_SimulateBenchmark(t *testing.T) {
	app := &BenchmarkApp{
		startTime:     time.Now(),
		lastFrameTime: time.Now(),
	}
	
	// Run simulated benchmark for a short duration
	startFrame := app.frameCount
	startTime := time.Now()
	
	// Simulate 30 frames (0.5 seconds at 60 FPS)
	for i := 0; i < 30; i++ {
		app.renderFrame()
		time.Sleep(time.Millisecond * 16) // ~60 FPS
	}
	
	duration := time.Since(startTime)
	
	// Verify frames were rendered
	if app.frameCount != startFrame+30 {
		t.Errorf("Expected %d frames, got %d", startFrame+30, app.frameCount)
	}
	
	// Verify reasonable duration (should be around 0.5 seconds)
	if duration < 400*time.Millisecond || duration > 800*time.Millisecond {
		t.Errorf("Unexpected duration: %v", duration)
	}
	
	// Verify rotation angle is within expected range
	expectedAngle := float32(30 * 0.01) // 30 frames * 0.01 rotation per frame
	if app.rotationAngle < expectedAngle*0.8 || app.rotationAngle > expectedAngle*1.2 {
		t.Errorf("Unexpected rotation angle: %f, expected around %f", app.rotationAngle, expectedAngle)
	}
}

func BenchmarkRenderFrame(b *testing.B) {
	app := &BenchmarkApp{
		startTime:     time.Now(),
		lastFrameTime: time.Now(),
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		app.renderFrame()
	}
}

func BenchmarkSimulateRenderingWork(b *testing.B) {
	app := &BenchmarkApp{
		rotationAngle: 1.0,
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		app.simulateRenderingWork()
	}
}