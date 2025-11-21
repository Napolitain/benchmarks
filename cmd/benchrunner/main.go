package main

import (
	"encoding/json"
	"fmt"
	"os"
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
	serverNames string
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

	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Run benchmarks",
		Long:  "Run benchmarks on all or selected servers",
		RunE:  runBenchmarks,
	}

	runCmd.Flags().IntVarP(&connections, "connections", "c", 100, "Number of connections")
	runCmd.Flags().IntVarP(&pipeline, "pipeline", "p", 1, "Pipeline factor")
	runCmd.Flags().IntVarP(&duration, "duration", "d", 10, "Duration in seconds")
	runCmd.Flags().StringVarP(&serverNames, "servers", "s", "", "Comma-separated server names (default: all)")

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

func runBenchmarks(cmd *cobra.Command, args []string) error {
	// Build binary first
	b := builder.New(baseDir)
	if err := b.Build(); err != nil {
		return fmt.Errorf("failed to build binary: %w", err)
	}

	// Get servers to benchmark
	var serversToRun []config.ServerConfig
	allServers := config.GetServers(baseDir)

	if serverNames == "" {
		serversToRun = allServers
	} else {
		names := strings.Split(serverNames, ",")
		for _, name := range names {
			name = strings.TrimSpace(name)
			if srv := config.GetServerByName(baseDir, name); srv != nil {
				serversToRun = append(serversToRun, *srv)
			} else {
				return fmt.Errorf("unknown server: %s", name)
			}
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
		result, err := runner.Run(srvCfg.Name, srvCfg.Port)
		if err != nil {
			fmt.Printf("WARNING: Benchmark failed: %v\n", err)
		}
		results = append(results, result)

		// Stop server
		srv.Stop()

		// Wait between benchmarks
		time.Sleep(2 * time.Second)
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
	fmt.Printf("%-20s %15s %15s\n", "Server", "Req/sec", "Status")
	fmt.Println(strings.Repeat("-", 80))

	for _, r := range results {
		status := "OK"
		reqPerSec := fmt.Sprintf("%.2f", r.ReqPerSec)
		if r.Error != "" {
			status = "FAILED"
			reqPerSec = "N/A"
		}
		fmt.Printf("%-20s %15s %15s\n", r.ServerName, reqPerSec, status)
	}
	fmt.Println(strings.Repeat("=", 80))
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
