package benchmark

import (
	"bufio"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type Result struct {
	ServerName   string    `json:"server_name"`
	ReqPerSec    float64   `json:"req_per_sec"`
	Connections  int       `json:"connections"`
	Pipeline     int       `json:"pipeline"`
	Duration     int       `json:"duration_seconds"`
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

func (r *Runner) Run(serverName string, port int) (*Result, error) {
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
	cmd := exec.Command(
		r.binaryPath,
		strconv.Itoa(r.connections),
		"localhost",
		strconv.Itoa(port),
		strconv.Itoa(r.pipeline),
	)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		result.Error = fmt.Sprintf("failed to create stdout pipe: %v", err)
		return result, err
	}

	if err := cmd.Start(); err != nil {
		result.Error = fmt.Sprintf("failed to start benchmark: %v", err)
		return result, err
	}

	// Parse output for "Req/sec: " lines
	var reqPerSecValues []float64
	scanner := bufio.NewScanner(stdout)
	
	done := make(chan bool)
	go func() {
		for scanner.Scan() {
			line := scanner.Text()
			fmt.Println("  " + line)
			
			// Parse "Req/sec: 123456.789"
			if strings.Contains(line, "Req/sec:") {
				parts := strings.Fields(line)
				if len(parts) >= 2 {
					if val, err := strconv.ParseFloat(parts[1], 64); err == nil {
						reqPerSecValues = append(reqPerSecValues, val)
					}
				}
			}
		}
		done <- true
	}()

	// Wait for duration + 2 seconds for warmup
	time.Sleep(time.Duration(r.duration+2) * time.Second)

	// Kill the benchmark process
	if cmd.Process != nil {
		cmd.Process.Kill()
	}

	<-done
	cmd.Wait()

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
	return result, nil
}
