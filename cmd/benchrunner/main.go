package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/benchmarks/internal/benchmark"
	"github.com/benchmarks/internal/builder"
	"github.com/benchmarks/internal/config"
	"github.com/benchmarks/internal/server"
	"github.com/spf13/cobra"
)

var (
	connections int
	pipeline    int
	duration    int
	baseDir     string
)

func main() {
	// Find base directory (repo root)
	wd, _ := os.Getwd()
	baseDir = findRepoRoot(wd)

	rootCmd := &cobra.Command{
		Use:   "benchrunner",
		Short: "HTTP benchmark orchestrator",
		Long:  "Orchestrates HTTP benchmarks across multiple server implementations",
	}

	// Run command with subcommands
	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Run benchmarks",
		Long:  "Run different types of benchmarks",
	}

	// Run server subcommand
	runServerCmd := &cobra.Command{
		Use:   "server [api-name]",
		Short: "Run HTTP server benchmarks",
		Long:  "Run benchmarks on all or selected HTTP servers",
		RunE:  runServerBenchmarks,
	}

	runServerCmd.Flags().IntVarP(&connections, "connections", "c", 100, "Number of connections")
	runServerCmd.Flags().IntVarP(&pipeline, "pipeline", "p", 1, "Pipeline factor")
	runServerCmd.Flags().IntVarP(&duration, "duration", "d", 10, "Duration in seconds")

	// Run startup subcommand
	runStartupCmd := &cobra.Command{
		Use:   "startup",
		Short: "Run startup time benchmarks",
		Long:  "Run startup time benchmarks for all languages using hyperfine",
		RunE:  runStartupBenchmarks,
	}

	runCmd.AddCommand(runServerCmd, runStartupCmd)

	buildCmd := &cobra.Command{
		Use:   "build",
		Short: "Build load test binary",
		Long:  "Build the http_load_test binary from uSockets",
		RunE:  buildBinary,
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List available servers",
		Long:  "List all available server implementations",
		Run:   listServers,
	}

	rootCmd.AddCommand(runCmd, buildCmd, listCmd)
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func findRepoRoot(start string) string {
	dir := start
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return start
		}
		dir = parent
	}
}

func buildBinary(cmd *cobra.Command, args []string) error {
	b := builder.New(baseDir)
	return b.Build()
}

func listServers(cmd *cobra.Command, args []string) {
	servers := config.GetServers(baseDir)
	fmt.Println("Available servers:")
	for _, s := range servers {
		fmt.Printf("  - %s\n", s.Name)
	}
}

func runServerBenchmarks(cmd *cobra.Command, args []string) error {
	// Parse server name from args
	var targetServer string
	if len(args) > 0 {
		targetServer = args[0]
	}
	// Build binary first
	b := builder.New(baseDir)
	if err := b.Build(); err != nil {
		return fmt.Errorf("failed to build binary: %w", err)
	}

	// Get servers to benchmark
	var serversToRun []config.ServerConfig
	allServers := config.GetServers(baseDir)

	if targetServer == "" {
		serversToRun = allServers
	} else {
		if srv := config.GetServerByName(baseDir, targetServer); srv != nil {
			serversToRun = append(serversToRun, *srv)
		} else {
			return fmt.Errorf("unknown server: %s", targetServer)
		}
	}

	if len(serversToRun) == 0 {
		return fmt.Errorf("no servers to benchmark")
	}

	// Run benchmarks
	runner := benchmark.NewRunner(b.GetBinaryPath(), connections, pipeline, duration)
	var results []*benchmark.Result

	for _, srvCfg := range serversToRun {
		srv := server.New(&srvCfg)
		
		// Start server
		if err := srv.Start(); err != nil {
			fmt.Printf("ERROR: %v\n", err)
			results = append(results, &benchmark.Result{
				ServerName: srvCfg.Name,
				Error:      err.Error(),
			})
			continue
		}

		// Run benchmark
		result, err := runner.Run(srvCfg.Name, srvCfg.Port, srv.GetPID())
		if err != nil {
			fmt.Printf("WARNING: Benchmark failed: %v\n", err)
		}
		results = append(results, result)

		// Stop server
		srv.Stop()

		// Wait between benchmarks (ensure port is fully released)
		time.Sleep(5 * time.Second)
	}

	// Print summary
	printSummary(results)

	// Save results
	if err := saveResults(results); err != nil {
		fmt.Printf("WARNING: Failed to save results: %v\n", err)
	}

	return nil
}

func printSummary(results []*benchmark.Result) {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("BENCHMARK SUMMARY")
	fmt.Println(strings.Repeat("=", 80))
	fmt.Printf("%-20s %15s %15s %15s\n", "Server", "Req/sec", "Memory (MB)", "Status")
	fmt.Println(strings.Repeat("-", 80))

	for _, r := range results {
		status := "OK"
		reqPerSec := fmt.Sprintf("%.2f", r.ReqPerSec)
		memory := fmt.Sprintf("%.2f", r.MemoryMB)
		if r.Error != "" {
			status = "FAILED"
			reqPerSec = "N/A"
			memory = "N/A"
		}
		fmt.Printf("%-20s %15s %15s %15s\n", r.ServerName, reqPerSec, memory, status)
	}
	fmt.Println(strings.Repeat("=", 80))
}

func runStartupBenchmarks(cmd *cobra.Command, args []string) error {
	startupDir := filepath.Join(baseDir, "startup")
	
	// Check if hyperfine is installed
	if _, err := exec.LookPath("hyperfine"); err != nil {
		return fmt.Errorf("hyperfine not found. Install with: sudo apt install hyperfine")
	}
	
	fmt.Println("Running startup time benchmarks with hyperfine...")
	fmt.Println(strings.Repeat("=", 80))
	
	// Define benchmark commands
	memoryDir := filepath.Join(startupDir, "memory")
	benchmarks := []struct {
		name string
		cmd  string
	}{
		{"Go Bubblesort", "cd " + filepath.Join(startupDir, "compute") + " && go run bubblesort.go"},
		{"Node Bubblesort", filepath.Join(startupDir, "compute", "bubblesort.js")},
		{"Python Bubblesort", filepath.Join(startupDir, "compute", "bubblesort.py")},
		{"Go Rectangle", "cd " + memoryDir + " && go run rectangle.go test_rectangle.yaml"},
		{"Node Rectangle", filepath.Join(memoryDir, "rectangle.js") + " " + filepath.Join(memoryDir, "test_rectangle.yaml")},
		{"Python Rectangle", filepath.Join(memoryDir, "rectangle.py") + " " + filepath.Join(memoryDir, "test_rectangle.yaml")},
	}
	
	// Build hyperfine command
	var cmdArgs []string
	cmdArgs = append(cmdArgs, "--warmup", "3", "--runs", "10")
	
	for _, b := range benchmarks {
		cmdArgs = append(cmdArgs, "--command-name", b.name, b.cmd)
	}
	
	// Run hyperfine
	hyperfineCmd := exec.Command("hyperfine", cmdArgs...)
	hyperfineCmd.Stdout = os.Stdout
	hyperfineCmd.Stderr = os.Stderr
	hyperfineCmd.Dir = baseDir
	
	if err := hyperfineCmd.Run(); err != nil {
		return fmt.Errorf("hyperfine failed: %w", err)
	}
	
	fmt.Println(strings.Repeat("=", 80))
	return nil
}

func saveResults(results []*benchmark.Result) error {
	resultsDir := filepath.Join(baseDir, "results")
	os.MkdirAll(resultsDir, 0755)

	filename := fmt.Sprintf("benchmark_%s.json", time.Now().Format("20060102_150405"))
	filepath := filepath.Join(resultsDir, filename)

	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(filepath, data, 0644); err != nil {
		return err
	}

	fmt.Printf("\nResults saved to: %s\n", filepath)
	return nil
}
