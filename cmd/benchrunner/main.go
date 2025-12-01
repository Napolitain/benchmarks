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

	runCmd.AddCommand(runServerCmd)

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

	// Run compute subcommand
	runComputeCmd := &cobra.Command{
		Use:   "compute [language]",
		Short: "Run compute (bubblesort) benchmarks",
		Long: `Compile and benchmark bubblesort programs in various languages using poop (or hyperfine as fallback).

Modes:
  compile    - Benchmark compilation time only (cold builds)
  full-cold  - Benchmark compilation + execution (cold builds, no cache)
  full-hot   - Benchmark compilation + execution (hot builds, cache allowed)
  exec       - Benchmark execution time only (pre-compiled)`,
		RunE: runComputeBenchmarks,
	}
	runComputeCmd.Flags().IntVarP(&warmup, "warmup", "w", 3, "Number of warmup runs")
	runComputeCmd.Flags().IntVarP(&runs, "runs", "r", 10, "Number of benchmark runs")
	runComputeCmd.Flags().StringVarP(&benchMode, "mode", "m", "exec", "Benchmark mode: compile, full-cold, full-hot, exec")

	runCmd.AddCommand(runComputeCmd)

	// Run cli subcommand
	runCLICmd := &cobra.Command{
		Use:   "cli [language]",
		Short: "Run CLI (rectangle YAML parsing) benchmarks",
		Long: `Compile and benchmark rectangle YAML parsing programs in various languages using poop (or hyperfine as fallback).

Modes:
  compile    - Benchmark compilation time only (cold builds)
  full-cold  - Benchmark compilation + execution (cold builds, no cache)
  full-hot   - Benchmark compilation + execution (hot builds, cache allowed)
  exec       - Benchmark execution time only (pre-compiled)`,
		RunE: runCLIBenchmarks,
	}
	runCLICmd.Flags().IntVarP(&warmup, "warmup", "w", 3, "Number of warmup runs")
	runCLICmd.Flags().IntVarP(&runs, "runs", "r", 10, "Number of benchmark runs")
	runCLICmd.Flags().StringVarP(&benchMode, "mode", "m", "exec", "Benchmark mode: compile, full-cold, full-hot, exec")

	runCmd.AddCommand(runCLICmd)

	// Run ffi subcommand
	runFFICmd := &cobra.Command{
		Use:   "ffi [language]",
		Short: "Run FFI benchmarks (fast_sum and slow_compute)",
		Long: `Compile and benchmark FFI programs in various languages using poop (or hyperfine as fallback).

Runs two sub-benchmarks sequentially:
  - fast_sum: measures FFI call overhead (1M calls of a simple sum function)
  - slow_compute: measures compute-heavy FFI (100 calls with 1M iterations each)

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
			name:   "node-direct",
			dir:    filepath.Join(hwDir, "node"),
			runCmd: "node main.js",
		},
		{
			name:       "node-build",
			dir:        filepath.Join(hwDir, "node"),
			compileCmd: "npx esbuild main.js --bundle --minify --platform=node --outfile=main.min.js",
			runCmd:     "node main.min.js",
			binaryPath: "main.min.js",
			cleanCmd:   "rm -f main.min.js",
			cleanFiles: []string{"main.min.js"},
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
// COMPUTE (Bubblesort) Benchmarks
// ============================================================================

func getComputeLanguages(baseDir string) []helloworldLang {
	computeDir := filepath.Join(baseDir, "compute")
	return []helloworldLang{
		{
			name:       "go",
			dir:        filepath.Join(computeDir, "go"),
			compileCmd: "go build -ldflags=\"-s -w\" -o bubblesort bubblesort.go",
			runCmd:     "./bubblesort",
			binaryPath: "bubblesort",
			cleanCmd:   "rm -f bubblesort",
			cleanFiles: []string{"bubblesort"},
		},
		{
			name:       "rust",
			dir:        filepath.Join(computeDir, "rust"),
			compileCmd: "rustc -C opt-level=3 -C lto=fat -C target-cpu=native -C strip=symbols -o bubblesort bubblesort.rs",
			runCmd:     "./bubblesort",
			binaryPath: "bubblesort",
			cleanCmd:   "rm -f bubblesort",
			cleanFiles: []string{"bubblesort"},
		},
		{
			name:       "zig",
			dir:        filepath.Join(computeDir, "zig"),
			compileCmd: "zig build-exe -OReleaseFast -fstrip -femit-bin=bubblesort bubblesort.zig",
			runCmd:     "./bubblesort",
			binaryPath: "bubblesort",
			cleanCmd:   "rm -f bubblesort bubblesort.o",
			cleanFiles: []string{"bubblesort", "bubblesort.o"},
		},
		{
			name:   "node-direct",
			dir:    filepath.Join(computeDir, "node"),
			runCmd: "node bubblesort.js",
		},
		{
			name:       "node-build",
			dir:        filepath.Join(computeDir, "node"),
			compileCmd: "npx esbuild bubblesort.js --bundle --minify --platform=node --outfile=bubblesort.min.js",
			runCmd:     "node bubblesort.min.js",
			binaryPath: "bubblesort.min.js",
			cleanCmd:   "rm -f bubblesort.min.js",
			cleanFiles: []string{"bubblesort.min.js"},
		},
		{
			name:   "python",
			dir:    filepath.Join(computeDir, "python"),
			runCmd: "python3 bubblesort.py",
		},
	}
}

func runComputeBenchmarks(cmd *cobra.Command, args []string) error {
	return runGenericBenchmarks("compute", getComputeLanguages(baseDir), args)
}

// ============================================================================
// CLI (Rectangle YAML parsing) Benchmarks
// ============================================================================

func getCLILanguages(baseDir string) []helloworldLang {
	cliDir := filepath.Join(baseDir, "cli")
	yamlFile := "../test_rectangle.yaml"
	return []helloworldLang{
		{
			name:       "cpp",
			dir:        filepath.Join(cliDir, "cpp"),
			compileCmd: "cmake -B build -DCMAKE_BUILD_TYPE=Release && cmake --build build -j$(nproc)",
			runCmd:     "./build/rectangle " + yamlFile,
			binaryPath: "build/rectangle",
			cleanCmd:   "rm -rf build",
			cleanFiles: []string{"build"},
		},
		{
			name:       "go",
			dir:        filepath.Join(cliDir, "go"),
			compileCmd: "go build -ldflags=\"-s -w\" -o rectangle rectangle.go",
			runCmd:     "./rectangle " + yamlFile,
			binaryPath: "rectangle",
			cleanCmd:   "rm -f rectangle",
			cleanFiles: []string{"rectangle"},
		},
		{
			name:       "rust",
			dir:        filepath.Join(cliDir, "rust"),
			compileCmd: "cargo build --release",
			runCmd:     "./target/release/rectangle " + yamlFile,
			binaryPath: "target/release/rectangle",
			cleanCmd:   "cargo clean",
			cleanFiles: []string{"target"},
		},
		{
			name:       "zig",
			dir:        filepath.Join(cliDir, "zig"),
			compileCmd: "zig build-exe -OReleaseFast -fstrip -femit-bin=rectangle rectangle.zig",
			runCmd:     "./rectangle " + yamlFile,
			binaryPath: "rectangle",
			cleanCmd:   "rm -f rectangle rectangle.o",
			cleanFiles: []string{"rectangle", "rectangle.o"},
		},
		{
			name:   "node-direct",
			dir:    filepath.Join(cliDir, "node"),
			runCmd: "node rectangle.js " + yamlFile,
		},
		{
			name:       "node-build",
			dir:        filepath.Join(cliDir, "node"),
			compileCmd: "npx esbuild rectangle.js --bundle --minify --platform=node --outfile=rectangle.min.js",
			runCmd:     "node rectangle.min.js " + yamlFile,
			binaryPath: "rectangle.min.js",
			cleanCmd:   "rm -f rectangle.min.js",
			cleanFiles: []string{"rectangle.min.js"},
		},
		{
			name:   "python",
			dir:    filepath.Join(cliDir, "python"),
			runCmd: "python3 rectangle.py " + yamlFile,
		},
	}
}

func runCLIBenchmarks(cmd *cobra.Command, args []string) error {
	return runGenericBenchmarks("cli", getCLILanguages(baseDir), args)
}

// ============================================================================
// FFI Benchmarks
// ============================================================================

func getFFIFastSumLanguages(baseDir string) []helloworldLang {
	ffiDir := filepath.Join(baseDir, "ffi", "fast_sum")
	return []helloworldLang{
		{
			name:       "cpp",
			dir:        filepath.Join(ffiDir, "cpp"),
			compileCmd: "g++ -O3 -flto -march=native -DNDEBUG -s -o main main.cpp ../hotpath.cpp",
			runCmd:     "./main",
			binaryPath: "main",
			cleanCmd:   "rm -f main",
			cleanFiles: []string{"main"},
		},
		{
			name:       "rust",
			dir:        filepath.Join(ffiDir, "rust"),
			compileCmd: "rustc -C opt-level=3 -C lto=fat -C target-cpu=native -C strip=symbols -o main main.rs",
			runCmd:     "./main",
			binaryPath: "main",
			cleanCmd:   "rm -f main",
			cleanFiles: []string{"main"},
		},
		{
			name:       "zig",
			dir:        filepath.Join(ffiDir, "zig"),
			compileCmd: "zig build-exe -OReleaseFast -fstrip -femit-bin=main main.zig",
			runCmd:     "./main",
			binaryPath: "main",
			cleanCmd:   "rm -f main main.o",
			cleanFiles: []string{"main", "main.o"},
		},
		{
			name:   "python",
			dir:    filepath.Join(ffiDir, "python"),
			runCmd: "python3 main.py",
		},
	}
}

func getFFISlowComputeLanguages(baseDir string) []helloworldLang {
	ffiDir := filepath.Join(baseDir, "ffi", "slow_compute")
	return []helloworldLang{
		{
			name:       "cpp",
			dir:        filepath.Join(ffiDir, "cpp"),
			compileCmd: "g++ -O3 -flto -march=native -DNDEBUG -s -o main main.cpp ../hotpath.cpp",
			runCmd:     "./main",
			binaryPath: "main",
			cleanCmd:   "rm -f main",
			cleanFiles: []string{"main"},
		},
		{
			name:       "rust",
			dir:        filepath.Join(ffiDir, "rust"),
			compileCmd: "rustc -C opt-level=3 -C lto=fat -C target-cpu=native -C strip=symbols -o main main.rs",
			runCmd:     "./main",
			binaryPath: "main",
			cleanCmd:   "rm -f main",
			cleanFiles: []string{"main"},
		},
		{
			name:       "zig",
			dir:        filepath.Join(ffiDir, "zig"),
			compileCmd: "zig build-exe -OReleaseFast -fstrip -femit-bin=main main.zig",
			runCmd:     "./main",
			binaryPath: "main",
			cleanCmd:   "rm -f main main.o",
			cleanFiles: []string{"main", "main.o"},
		},
		{
			name:   "python",
			dir:    filepath.Join(ffiDir, "python"),
			runCmd: "python3 main.py",
		},
	}
}

func runFFIBenchmarks(cmd *cobra.Command, args []string) error {
	fmt.Println("Running FFI benchmarks (two sub-benchmarks)")
	fmt.Println(strings.Repeat("=", 80))

	// Run fast_sum benchmark
	fmt.Println("\n[1/2] fast_sum - FFI call overhead benchmark")
	fmt.Println(strings.Repeat("-", 80))
	if err := runGenericBenchmarks("ffi/fast_sum", getFFIFastSumLanguages(baseDir), args); err != nil {
		return err
	}

	// Run slow_compute benchmark
	fmt.Println("\n[2/2] slow_compute - compute-heavy FFI benchmark")
	fmt.Println(strings.Repeat("-", 80))
	if err := runGenericBenchmarks("ffi/slow_compute", getFFISlowComputeLanguages(baseDir), args); err != nil {
		return err
	}

	return nil
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
			// Skip interpreted languages only for compile mode
			if benchMode == "compile" && lang.compileCmd == "" {
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
				if lang.compileCmd == "" {
					// Interpreted language - just run
					benchCmd = fmt.Sprintf("cd %s && %s", lang.dir, lang.runCmd)
				} else {
					benchCmd = fmt.Sprintf("cd %s && %s && %s && %s", lang.dir, lang.cleanCmd, lang.compileCmd, lang.runCmd)
				}
			case "full-hot":
				if lang.compileCmd == "" {
					// Interpreted language - just run
					benchCmd = fmt.Sprintf("cd %s && %s", lang.dir, lang.runCmd)
				} else {
					benchCmd = fmt.Sprintf("cd %s && %s && %s", lang.dir, lang.compileCmd, lang.runCmd)
				}
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
				if lang.compileCmd == "" {
					// Interpreted language - just run
					benchCmd = fmt.Sprintf("cd %s && %s", lang.dir, lang.runCmd)
				} else {
					benchCmd = fmt.Sprintf("cd %s && %s && %s", lang.dir, lang.compileCmd, lang.runCmd)
					prepareCmd = fmt.Sprintf("cd %s && %s", lang.dir, lang.cleanCmd)
				}

			case "full-hot":
				if lang.compileCmd == "" {
					// Interpreted language - just run
					benchCmd = fmt.Sprintf("cd %s && %s", lang.dir, lang.runCmd)
				} else {
					benchCmd = fmt.Sprintf("cd %s && %s && %s", lang.dir, lang.compileCmd, lang.runCmd)
				}

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
