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
	warmup      int
	runs        int
	benchMode   string
)

// getBenchmarkTool returns "poop" if available, otherwise "hyperfine"
func getBenchmarkTool() (string, error) {
	if _, err := exec.LookPath("poop"); err == nil {
		return "poop", nil
	}
	if _, err := exec.LookPath("hyperfine"); err == nil {
		return "hyperfine", nil
	}
	return "", fmt.Errorf("neither poop nor hyperfine found. Install poop or hyperfine")
}

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
		Long:  "Run startup time benchmarks for all languages using poop (or hyperfine as fallback)",
		RunE:  runStartupBenchmarks,
	}

	runCmd.AddCommand(runServerCmd, runStartupCmd)

	// Run helloworld subcommand
	runHelloworldCmd := &cobra.Command{
		Use:   "helloworld [language]",
		Short: "Run helloworld benchmarks",
		Long: `Compile and benchmark helloworld programs in various languages using poop (or hyperfine as fallback).

Modes:
  compile    - Benchmark compilation time only (cold builds)
  full-cold  - Benchmark compilation + execution (cold builds, no cache)
  full-hot   - Benchmark compilation + execution (hot builds, cache allowed)
  exec       - Benchmark execution time only (pre-compiled)`,
		RunE: runHelloworldBenchmarks,
	}
	runHelloworldCmd.Flags().IntVarP(&warmup, "warmup", "w", 3, "Number of warmup runs")
	runHelloworldCmd.Flags().IntVarP(&runs, "runs", "r", 10, "Number of benchmark runs")
	runHelloworldCmd.Flags().StringVarP(&benchMode, "mode", "m", "exec", "Benchmark mode: compile, full-cold, full-hot, exec")

	runCmd.AddCommand(runHelloworldCmd)

	// Run startup-compute subcommand
	runStartupComputeCmd := &cobra.Command{
		Use:   "compute [language]",
		Short: "Run compute (bubblesort) benchmarks",
		Long: `Compile and benchmark bubblesort programs in various languages using poop (or hyperfine as fallback).

Modes:
  compile    - Benchmark compilation time only (cold builds)
  full-cold  - Benchmark compilation + execution (cold builds, no cache)
  full-hot   - Benchmark compilation + execution (hot builds, cache allowed)
  exec       - Benchmark execution time only (pre-compiled)`,
		RunE: runStartupComputeBenchmarks,
	}
	runStartupComputeCmd.Flags().IntVarP(&warmup, "warmup", "w", 3, "Number of warmup runs")
	runStartupComputeCmd.Flags().IntVarP(&runs, "runs", "r", 10, "Number of benchmark runs")
	runStartupComputeCmd.Flags().StringVarP(&benchMode, "mode", "m", "exec", "Benchmark mode: compile, full-cold, full-hot, exec")

	runCmd.AddCommand(runStartupComputeCmd)

	// Run startup-memory subcommand
	runStartupMemoryCmd := &cobra.Command{
		Use:   "cli [language]",
		Short: "Run CLI (rectangle YAML parsing) benchmarks",
		Long: `Compile and benchmark rectangle YAML parsing programs in various languages using poop (or hyperfine as fallback).

Modes:
  compile    - Benchmark compilation time only (cold builds)
  full-cold  - Benchmark compilation + execution (cold builds, no cache)
  full-hot   - Benchmark compilation + execution (hot builds, cache allowed)
  exec       - Benchmark execution time only (pre-compiled)`,
		RunE: runStartupMemoryBenchmarks,
	}
	runStartupMemoryCmd.Flags().IntVarP(&warmup, "warmup", "w", 3, "Number of warmup runs")
	runStartupMemoryCmd.Flags().IntVarP(&runs, "runs", "r", 10, "Number of benchmark runs")
	runStartupMemoryCmd.Flags().StringVarP(&benchMode, "mode", "m", "exec", "Benchmark mode: compile, full-cold, full-hot, exec")

	runCmd.AddCommand(runStartupMemoryCmd)

	// Run ffi subcommand
	runFFICmd := &cobra.Command{
		Use:   "ffi [language]",
		Short: "Run FFI benchmarks",
		Long: `Compile and benchmark FFI programs in various languages using poop (or hyperfine as fallback).

Modes:
  compile    - Benchmark compilation time only (cold builds)
  full-cold  - Benchmark compilation + execution (cold builds, no cache)
  full-hot   - Benchmark compilation + execution (hot builds, cache allowed)
  exec       - Benchmark execution time only (pre-compiled)`,
		RunE: runFFIBenchmarks,
	}
	runFFICmd.Flags().IntVarP(&warmup, "warmup", "w", 3, "Number of warmup runs")
	runFFICmd.Flags().IntVarP(&runs, "runs", "r", 10, "Number of benchmark runs")
	runFFICmd.Flags().StringVarP(&benchMode, "mode", "m", "exec", "Benchmark mode: compile, full-cold, full-hot, exec")

	runCmd.AddCommand(runFFICmd)

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
	// Check for benchmark tool
	benchTool, err := getBenchmarkTool()
	if err != nil {
		return err
	}
	
	fmt.Printf("Running startup time benchmarks with %s...\n", benchTool)
	fmt.Println(strings.Repeat("=", 80))
	
	// Define benchmark commands
	computeDir := filepath.Join(baseDir, "compute")
	cliDir := filepath.Join(baseDir, "cli")
	benchmarks := []struct {
		name string
		cmd  string
	}{
		{"Go Bubblesort", "cd " + computeDir + " && go run bubblesort.go"},
		{"Node Bubblesort", filepath.Join(computeDir, "bubblesort.js")},
		{"Python Bubblesort", filepath.Join(computeDir, "bubblesort.py")},
		{"Go Rectangle", "cd " + cliDir + " && go run rectangle.go test_rectangle.yaml"},
		{"Node Rectangle", filepath.Join(cliDir, "rectangle.js") + " " + filepath.Join(cliDir, "test_rectangle.yaml")},
		{"Python Rectangle", filepath.Join(cliDir, "rectangle.py") + " " + filepath.Join(cliDir, "test_rectangle.yaml")},
	}
	
	// Build command args based on tool
	var cmdArgs []string
	if benchTool == "poop" {
		for _, b := range benchmarks {
			cmdArgs = append(cmdArgs, b.cmd)
		}
	} else {
		cmdArgs = append(cmdArgs, "--warmup", "3", "--runs", "10")
		for _, b := range benchmarks {
			cmdArgs = append(cmdArgs, "--command-name", b.name, b.cmd)
		}
	}
	
	// Run benchmark tool
	benchCmd := exec.Command(benchTool, cmdArgs...)
	benchCmd.Stdout = os.Stdout
	benchCmd.Stderr = os.Stderr
	benchCmd.Dir = baseDir
	
	if err := benchCmd.Run(); err != nil {
		return fmt.Errorf("%s failed: %w", benchTool, err)
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

// helloworldLang defines a language's build and run configuration
type helloworldLang struct {
	name       string
	dir        string
	compileCmd string   // command to compile
	runCmd     string   // command to run the binary
	binaryPath string   // path to the compiled binary (relative to dir)
	cleanCmd   string   // command to clean build artifacts
	cleanFiles []string // files/dirs to remove for cold builds
}

func getHelloworldLanguages(baseDir string) []helloworldLang {
	hwDir := filepath.Join(baseDir, "helloworld")
	return []helloworldLang{
		// C variants
		{
			name:       "c-cmake",
			dir:        filepath.Join(hwDir, "c"),
			compileCmd: "cmake -B build -DCMAKE_BUILD_TYPE=Release && cmake --build build",
			runCmd:     "./build/hello",
			binaryPath: "build/hello",
			cleanCmd:   "rm -rf build",
			cleanFiles: []string{"build"},
		},
		{
			name:       "c-direct",
			dir:        filepath.Join(hwDir, "c"),
			compileCmd: "gcc -O3 -flto -march=native -DNDEBUG -s -o hello main.c",
			runCmd:     "./hello",
			binaryPath: "hello",
			cleanCmd:   "rm -f hello",
			cleanFiles: []string{"hello"},
		},
		// C++ variants
		{
			name:       "cpp-cmake",
			dir:        filepath.Join(hwDir, "cpp"),
			compileCmd: "cmake -B build -DCMAKE_BUILD_TYPE=Release && cmake --build build",
			runCmd:     "./build/hello",
			binaryPath: "build/hello",
			cleanCmd:   "rm -rf build",
			cleanFiles: []string{"build"},
		},
		{
			name:       "cpp-direct",
			dir:        filepath.Join(hwDir, "cpp"),
			compileCmd: "g++ -O3 -flto -march=native -DNDEBUG -s -o hello main.cpp",
			runCmd:     "./hello",
			binaryPath: "hello",
			cleanCmd:   "rm -f hello",
			cleanFiles: []string{"hello"},
		},
		// Go (already direct)
		{
			name:       "go",
			dir:        filepath.Join(hwDir, "go"),
			compileCmd: "go build -ldflags=\"-s -w\" -o hello main.go",
			runCmd:     "./hello",
			binaryPath: "hello",
			cleanCmd:   "rm -f hello",
			cleanFiles: []string{"hello"},
		},
		// Rust variants
		{
			name:       "rust-cargo",
			dir:        filepath.Join(hwDir, "rust"),
			compileCmd: "cargo build --release",
			runCmd:     "./target/release/hello",
			binaryPath: "target/release/hello",
			cleanCmd:   "cargo clean",
			cleanFiles: []string{"target"},
		},
		{
			name:       "rust-direct",
			dir:        filepath.Join(hwDir, "rust"),
			compileCmd: "rustc -C opt-level=3 -C lto=fat -C target-cpu=native -C strip=symbols -o hello src/main.rs",
			runCmd:     "./hello",
			binaryPath: "hello",
			cleanCmd:   "rm -f hello",
			cleanFiles: []string{"hello"},
		},
		// Zig variants
		{
			name:       "zig-build",
			dir:        filepath.Join(hwDir, "zig"),
			compileCmd: "zig build -Doptimize=ReleaseFast",
			runCmd:     "./zig-out/bin/hello",
			binaryPath: "zig-out/bin/hello",
			cleanCmd:   "rm -rf zig-out .zig-cache",
			cleanFiles: []string{"zig-out", ".zig-cache"},
		},
		{
			name:       "zig-direct",
			dir:        filepath.Join(hwDir, "zig"),
			compileCmd: "zig build-exe -OReleaseFast -fstrip -femit-bin=hello main.zig",
			runCmd:     "./hello",
			binaryPath: "hello",
			cleanCmd:   "rm -f hello hello.o",
			cleanFiles: []string{"hello", "hello.o"},
		},
		// Interpreted languages
		{
			name:   "node",
			dir:    filepath.Join(hwDir, "node"),
			runCmd: "node main.js",
		},
		{
			name:   "python",
			dir:    filepath.Join(hwDir, "python"),
			runCmd: "python3 main.py",
		},
	}
}

func runHelloworldBenchmarks(cmd *cobra.Command, args []string) error {
	return runGenericBenchmarks("helloworld", getHelloworldLanguages(baseDir), args)
}

// ============================================================================
// STARTUP-COMPUTE (Bubblesort) Benchmarks
// ============================================================================

func getStartupComputeLanguages(baseDir string) []helloworldLang {
	dir := filepath.Join(baseDir, "compute")
	return []helloworldLang{
		{
			name:       "go",
			dir:        dir,
			compileCmd: "go build -ldflags=\"-s -w\" -o bubblesort_bin bubblesort.go",
			runCmd:     "./bubblesort_bin",
			binaryPath: "bubblesort_bin",
			cleanCmd:   "rm -f bubblesort_bin",
			cleanFiles: []string{"bubblesort_bin"},
		},
		{
			name:       "rust",
			dir:        dir,
			compileCmd: "rustc -C opt-level=3 -C lto=fat -C target-cpu=native -C strip=symbols -o bubblesort_bin bubblesort.rs",
			runCmd:     "./bubblesort_bin",
			binaryPath: "bubblesort_bin",
			cleanCmd:   "rm -f bubblesort_bin",
			cleanFiles: []string{"bubblesort_bin"},
		},
		{
			name:       "zig",
			dir:        dir,
			compileCmd: "zig build-exe -OReleaseFast -fstrip -femit-bin=bubblesort_bin bubblesort.zig",
			runCmd:     "./bubblesort_bin",
			binaryPath: "bubblesort_bin",
			cleanCmd:   "rm -f bubblesort_bin bubblesort_bin.o",
			cleanFiles: []string{"bubblesort_bin", "bubblesort_bin.o"},
		},
		{
			name:   "node",
			dir:    dir,
			runCmd: "node bubblesort.js",
		},
		{
			name:   "python",
			dir:    dir,
			runCmd: "python3 bubblesort.py",
		},
	}
}

func runStartupComputeBenchmarks(cmd *cobra.Command, args []string) error {
	return runGenericBenchmarks("startup-compute", getStartupComputeLanguages(baseDir), args)
}

// ============================================================================
// STARTUP-MEMORY (Rectangle YAML parsing) Benchmarks
// ============================================================================

func getStartupMemoryLanguages(baseDir string) []helloworldLang {
	dir := filepath.Join(baseDir, "cli")
	cppDir := filepath.Join(dir, "cpp")
	yamlFile := filepath.Join(dir, "test_rectangle.yaml")
	return []helloworldLang{
		{
			name:       "cpp",
			dir:        cppDir,
			compileCmd: "cmake -B build -DCMAKE_BUILD_TYPE=Release && cmake --build build -j$(nproc)",
			runCmd:     "./build/rectangle " + yamlFile,
			binaryPath: "build/rectangle",
			cleanCmd:   "rm -rf build",
			cleanFiles: []string{"build"},
		},
		{
			name:       "go",
			dir:        dir,
			compileCmd: "go build -ldflags=\"-s -w\" -o rectangle_bin rectangle.go",
			runCmd:     "./rectangle_bin " + yamlFile,
			binaryPath: "rectangle_bin",
			cleanCmd:   "rm -f rectangle_bin",
			cleanFiles: []string{"rectangle_bin"},
		},
		{
			name:       "rust",
			dir:        dir,
			compileCmd: "cargo build --release",
			runCmd:     "./target/release/rectangle " + yamlFile,
			binaryPath: "target/release/rectangle",
			cleanCmd:   "cargo clean",
			cleanFiles: []string{"target"},
		},
		{
			name:       "zig",
			dir:        dir,
			compileCmd: "zig build-exe -OReleaseFast -fstrip -femit-bin=rectangle_bin rectangle.zig",
			runCmd:     "./rectangle_bin " + yamlFile,
			binaryPath: "rectangle_bin",
			cleanCmd:   "rm -f rectangle_bin rectangle_bin.o",
			cleanFiles: []string{"rectangle_bin", "rectangle_bin.o"},
		},
		{
			name:   "node",
			dir:    dir,
			runCmd: "node rectangle.js " + yamlFile,
		},
		{
			name:   "python",
			dir:    dir,
			runCmd: "python3 rectangle.py " + yamlFile,
		},
	}
}

func runStartupMemoryBenchmarks(cmd *cobra.Command, args []string) error {
	return runGenericBenchmarks("startup-memory", getStartupMemoryLanguages(baseDir), args)
}

// ============================================================================
// FFI Benchmarks
// ============================================================================

func getFFILanguages(baseDir string) []helloworldLang {
	dir := filepath.Join(baseDir, "ffi")
	return []helloworldLang{
		{
			name:       "native",
			dir:        dir,
			compileCmd: "g++ -O3 -flto -march=native -DNDEBUG -s -o bench_native bench_native.cpp hotpath.cpp",
			runCmd:     "./bench_native",
			binaryPath: "bench_native",
			cleanCmd:   "rm -f bench_native",
			cleanFiles: []string{"bench_native"},
		},
		{
			name:       "go-cgo",
			dir:        dir,
			compileCmd: "go build -ldflags=\"-s -w\" -o bench_go_bin bench_go_cgo.go",
			runCmd:     "LD_LIBRARY_PATH=. ./bench_go_bin",
			binaryPath: "bench_go_bin",
			cleanCmd:   "rm -f bench_go_bin",
			cleanFiles: []string{"bench_go_bin"},
		},
		{
			name:       "rust",
			dir:        dir,
			compileCmd: "rustc -C opt-level=3 -C lto=fat -C target-cpu=native -C strip=symbols -o bench_rust_bin bench_rust.rs",
			runCmd:     "./bench_rust_bin",
			binaryPath: "bench_rust_bin",
			cleanCmd:   "rm -f bench_rust_bin",
			cleanFiles: []string{"bench_rust_bin"},
		},
		{
			name:       "zig",
			dir:        dir,
			compileCmd: "zig build-exe -OReleaseFast -fstrip -femit-bin=bench_zig_bin bench_zig.zig",
			runCmd:     "./bench_zig_bin",
			binaryPath: "bench_zig_bin",
			cleanCmd:   "rm -f bench_zig_bin bench_zig_bin.o",
			cleanFiles: []string{"bench_zig_bin", "bench_zig_bin.o"},
		},
		{
			name:   "python",
			dir:    dir,
			runCmd: "python3 bench_python.py",
		},
		{
			name:   "node",
			dir:    dir,
			runCmd: "npx ts-node bench_node.ts",
		},
	}
}

func runFFIBenchmarks(cmd *cobra.Command, args []string) error {
	return runGenericBenchmarks("ffi", getFFILanguages(baseDir), args)
}

// ============================================================================
// Generic Benchmark Runner
// ============================================================================

func runGenericBenchmarks(suiteName string, languages []helloworldLang, args []string) error {
	// Validate mode
	validModes := map[string]bool{"compile": true, "full-cold": true, "full-hot": true, "exec": true}
	if !validModes[benchMode] {
		return fmt.Errorf("invalid mode: %s (valid: compile, full-cold, full-hot, exec)", benchMode)
	}

	// Check for benchmark tool
	benchTool, err := getBenchmarkTool()
	if err != nil {
		return err
	}

	// Filter by language if specified
	var targetLang string
	if len(args) > 0 {
		targetLang = strings.ToLower(args[0])
	}

	var langsToRun []helloworldLang
	for _, lang := range languages {
		if targetLang == "" || lang.name == targetLang {
			// Skip interpreted languages for compile modes
			if (benchMode == "compile" || benchMode == "full-cold" || benchMode == "full-hot") && lang.compileCmd == "" {
				fmt.Printf("Skipping %s (interpreted, no compilation)\n", lang.name)
				continue
			}
			langsToRun = append(langsToRun, lang)
		}
	}

	if len(langsToRun) == 0 {
		return fmt.Errorf("no matching language found: %s", targetLang)
	}

	fmt.Printf("Running %s benchmarks [mode: %s]\n", suiteName, benchMode)
	fmt.Println(strings.Repeat("=", 80))

	// For exec mode, pre-compile all binaries first
	if benchMode == "exec" {
		fmt.Println("Pre-compiling binaries...")
		for _, lang := range langsToRun {
			if lang.compileCmd == "" {
				fmt.Printf("%-20s: interpreted (no build needed)\n", lang.name)
				continue
			}
			fmt.Printf("%-20s: compiling... ", lang.name)
			compileExec := exec.Command("sh", "-c", lang.compileCmd)
			compileExec.Dir = lang.dir
			output, err := compileExec.CombinedOutput()
			if err != nil {
				fmt.Printf("FAILED\n%s\n", string(output))
				return fmt.Errorf("failed to compile %s: %w", lang.name, err)
			}
			fmt.Println("OK")
		}
		fmt.Println(strings.Repeat("=", 80))
	}

	fmt.Printf("\nRunning benchmarks with %s...\n", benchTool)
	fmt.Println(strings.Repeat("=", 80))

	// Build command args based on tool
	var cmdArgs []string
	if benchTool == "poop" {
		for _, lang := range langsToRun {
			var benchCmd string
			switch benchMode {
			case "compile":
				benchCmd = fmt.Sprintf("cd %s && %s", lang.dir, lang.compileCmd)
			case "full-cold":
				benchCmd = fmt.Sprintf("cd %s && %s && %s && %s", lang.dir, lang.cleanCmd, lang.compileCmd, lang.runCmd)
			case "full-hot":
				benchCmd = fmt.Sprintf("cd %s && %s && %s", lang.dir, lang.compileCmd, lang.runCmd)
			case "exec":
				benchCmd = fmt.Sprintf("cd %s && %s", lang.dir, lang.runCmd)
			}
			cmdArgs = append(cmdArgs, benchCmd)
		}
	} else {
		cmdArgs = append(cmdArgs, "--warmup", fmt.Sprintf("%d", warmup), "--runs", fmt.Sprintf("%d", runs))

		for _, lang := range langsToRun {
			var benchCmd, prepareCmd string

			switch benchMode {
			case "compile":
				benchCmd = fmt.Sprintf("cd %s && %s", lang.dir, lang.compileCmd)
				prepareCmd = fmt.Sprintf("cd %s && %s", lang.dir, lang.cleanCmd)

			case "full-cold":
				benchCmd = fmt.Sprintf("cd %s && %s && %s", lang.dir, lang.compileCmd, lang.runCmd)
				prepareCmd = fmt.Sprintf("cd %s && %s", lang.dir, lang.cleanCmd)

			case "full-hot":
				benchCmd = fmt.Sprintf("cd %s && %s && %s", lang.dir, lang.compileCmd, lang.runCmd)

			case "exec":
				benchCmd = fmt.Sprintf("cd %s && %s", lang.dir, lang.runCmd)
			}

			if prepareCmd != "" {
				cmdArgs = append(cmdArgs, "--prepare", prepareCmd)
			}
			cmdArgs = append(cmdArgs, "--command-name", lang.name, benchCmd)
		}
	}

	// Run benchmark tool
	benchExec := exec.Command(benchTool, cmdArgs...)
	benchExec.Stdout = os.Stdout
	benchExec.Stderr = os.Stderr
	benchExec.Dir = baseDir

	if err := benchExec.Run(); err != nil {
		return fmt.Errorf("%s failed: %w", benchTool, err)
	}

	fmt.Println(strings.Repeat("=", 80))

	// Report binary sizes for compile modes
	if benchMode == "compile" || benchMode == "full-cold" || benchMode == "full-hot" {
		fmt.Println("\nBinary sizes:")
		fmt.Println(strings.Repeat("-", 40))
		fmt.Printf("%-20s %10s\n", "Language", "Size")
		fmt.Println(strings.Repeat("-", 40))

		for _, lang := range langsToRun {
			if lang.binaryPath == "" {
				continue
			}
			binaryFullPath := filepath.Join(lang.dir, lang.binaryPath)
			if info, err := os.Stat(binaryFullPath); err == nil {
				size := info.Size()
				var sizeStr string
				if size >= 1024*1024 {
					sizeStr = fmt.Sprintf("%.2f MB", float64(size)/(1024*1024))
				} else if size >= 1024 {
					sizeStr = fmt.Sprintf("%.2f KB", float64(size)/1024)
				} else {
					sizeStr = fmt.Sprintf("%d B", size)
				}
				fmt.Printf("%-20s %10s\n", lang.name, sizeStr)
			} else {
				fmt.Printf("%-20s %10s\n", lang.name, "N/A")
			}
		}
		fmt.Println(strings.Repeat("-", 40))
	}

	// Clean up build artifacts
	fmt.Println("\nCleaning up build artifacts...")
	for _, lang := range langsToRun {
		if lang.cleanCmd != "" {
			cleanExec := exec.Command("sh", "-c", lang.cleanCmd)
			cleanExec.Dir = lang.dir
			cleanExec.Run()
		}
	}

	return nil
}
