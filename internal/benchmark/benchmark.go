package benchmark

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"
)

type Result struct {
	ServerName   string    `json:"server_name"`
	ReqPerSec    float64   `json:"req_per_sec"`
	Connections  int       `json:"connections"`
	Pipeline     int       `json:"pipeline"`
	Duration     int       `json:"duration_seconds"`
	MemoryMB     float64   `json:"memory_mb"`
	Timestamp    time.Time `json:"timestamp"`
	Error        string    `json:"error,omitempty"`
}

type Runner struct {
	binaryPath  string
	connections int
	pipeline    int
	duration    int
}

func NewRunner(binaryPath string, connections, pipeline, duration int) *Runner {
	return &Runner{
		binaryPath:  binaryPath,
		connections: connections,
		pipeline:    pipeline,
		duration:    duration,
	}
}

func (r *Runner) Run(serverName string, port int, serverPID int) (*Result, error) {
	fmt.Printf("\nRunning benchmark for %s...\n", serverName)
	fmt.Printf("  Connections: %d, Pipeline: %d, Duration: %ds\n", r.connections, r.pipeline, r.duration)

	result := &Result{
		ServerName:  serverName,
		Connections: r.connections,
		Pipeline:    r.pipeline,
		Duration:    r.duration,
		Timestamp:   time.Now(),
	}

	// Run: http_load_test <connections> <host> <port> [pipeline]
	// Use stdbuf to unbuffer output
	cmd := exec.Command(
		"stdbuf", "-oL", "-eL",
		r.binaryPath,
		strconv.Itoa(r.connections),
		"localhost",
		strconv.Itoa(port),
		strconv.Itoa(r.pipeline),
	)
	
	var stdBuffer bytes.Buffer
	mw := io.MultiWriter(os.Stdout, &stdBuffer)
	cmd.Stdout = mw
	cmd.Stderr = mw
	
	fmt.Printf("  Running benchmark for %d seconds...\n", r.duration)
	
	// Set process group so we can kill all children
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	
	// Start the command
	if err := cmd.Start(); err != nil {
		result.Error = fmt.Sprintf("failed to start: %v", err)
		return result, err
	}
	
	fmt.Printf("  Benchmark process started (PID %d), waiting %d seconds...\n", cmd.Process.Pid, r.duration+5)
	
	// Wait for duration + extra seconds for warmup
	time.Sleep(time.Duration(r.duration+5) * time.Second)
	
	// Get memory usage before killing
	if serverPID > 0 {
		result.MemoryMB = getMemoryUsage(serverPID)
	}
	
	// Kill the entire process group
	if cmd.Process != nil {
		fmt.Printf("  Killing benchmark process group (PID %d)...\n", cmd.Process.Pid)
		// Kill the process group (negative PID)
		syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
	}
	cmd.Wait()
	fmt.Printf("  Process terminated.\n")
	
	content := stdBuffer.Bytes()
	fmt.Printf("  Captured %d bytes\n", len(content))
	
	var reqPerSecValues []float64
	scanner := bufio.NewScanner(bytes.NewReader(content))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "Req/sec:") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				if val, err := strconv.ParseFloat(parts[1], 64); err == nil {
					reqPerSecValues = append(reqPerSecValues, val)
				}
			}
		}
	}

	// Calculate average req/sec (skip first value as warmup)
	if len(reqPerSecValues) > 1 {
		sum := 0.0
		for i := 1; i < len(reqPerSecValues); i++ {
			sum += reqPerSecValues[i]
		}
		result.ReqPerSec = sum / float64(len(reqPerSecValues)-1)
	} else if len(reqPerSecValues) == 1 {
		result.ReqPerSec = reqPerSecValues[0]
	} else {
		result.Error = "no benchmark results captured"
		return result, fmt.Errorf("no results")
	}

	fmt.Printf("  Average Req/sec: %.2f\n", result.ReqPerSec)
	fmt.Printf("  Memory usage: %.2f MB\n", result.MemoryMB)
	return result, nil
}

func getMemoryUsage(pid int) float64 {
	// Use ps to get RSS (Resident Set Size) in KB
	psCmd := exec.Command("ps", "-p", strconv.Itoa(pid), "-o", "rss=")
	output, err := psCmd.Output()
	if err != nil {
		return 0
	}
	
	rssKB := strings.TrimSpace(string(output))
	kb, err := strconv.ParseFloat(rssKB, 64)
	if err != nil {
		return 0
	}
	
	return kb / 1024.0 // Convert KB to MB
}
