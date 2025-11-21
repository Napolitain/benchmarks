package builder

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type Builder struct {
	uWebSocketsPath string
	uSocketsPath    string
	outputBinary    string
	uWebSocketsBin  string
}

func New(baseDir string) *Builder {
	uWebSocketsPath := filepath.Join(baseDir, "api", "uWebSockets")
	uSocketsPath := filepath.Join(uWebSocketsPath, "uSockets")
	outputBinary := filepath.Join(baseDir, "bin", "http_load_test")
	uWebSocketsBin := filepath.Join(baseDir, "bin", "HelloWorldBenchmark")
	
	return &Builder{
		uWebSocketsPath: uWebSocketsPath,
		uSocketsPath:    uSocketsPath,
		outputBinary:    outputBinary,
		uWebSocketsBin:  uWebSocketsBin,
	}
}

func (b *Builder) Build() error {
	// Create bin directory
	binDir := filepath.Join(filepath.Dir(b.outputBinary))
	if err := os.MkdirAll(binDir, 0755); err != nil {
		return fmt.Errorf("failed to create bin directory: %w", err)
	}

	// Build uSockets library first
	if err := b.buildUSockets(); err != nil {
		return err
	}

	// Build http_load_test tool
	if err := b.buildLoadTest(); err != nil {
		return err
	}

	// Build uWebSockets HelloWorldBenchmark server
	if err := b.buildUWebSocketsServer(); err != nil {
		return err
	}

	return nil
}

func (b *Builder) buildUSockets() error {
	fmt.Println("Building uSockets library...")

	// Check if already built
	uSocketsLib := filepath.Join(b.uSocketsPath, "uSockets.a")
	if _, err := os.Stat(uSocketsLib); err == nil {
		fmt.Println("uSockets.a already exists, skipping build")
		return nil
	}

	makeCmd := exec.Command("make", "default")
	makeCmd.Dir = b.uSocketsPath
	makeCmd.Env = append(os.Environ(), "WITH_OPENSSL=0")
	makeCmd.Stdout = os.Stdout
	makeCmd.Stderr = os.Stderr
	
	if err := makeCmd.Run(); err != nil {
		return fmt.Errorf("failed to build uSockets: %w", err)
	}

	fmt.Println("uSockets library built successfully")
	return nil
}

func (b *Builder) buildLoadTest() error {
	// Check if already built
	if _, err := os.Stat(b.outputBinary); err == nil {
		fmt.Println("http_load_test already exists, skipping build")
		return nil
	}

	fmt.Println("Building http_load_test...")
	
	sourceFile := filepath.Join(b.uSocketsPath, "examples", "http_load_test.c")
	
	compileCmd := exec.Command("gcc",
		"-std=c11",
		"-O3",
		"-DLIBUS_NO_SSL",
		"-I"+filepath.Join(b.uSocketsPath, "src"),
		"-o", b.outputBinary,
		sourceFile,
		filepath.Join(b.uSocketsPath, "uSockets.a"),
	)
	
	compileCmd.Dir = b.uSocketsPath
	compileCmd.Stdout = os.Stdout
	compileCmd.Stderr = os.Stderr
	
	if err := compileCmd.Run(); err != nil {
		return fmt.Errorf("failed to compile http_load_test: %w", err)
	}

	fmt.Printf("Successfully built: %s\n", b.outputBinary)
	return nil
}

func (b *Builder) buildUWebSocketsServer() error {
	// Check if already built
	if _, err := os.Stat(b.uWebSocketsBin); err == nil {
		fmt.Println("HelloWorldBenchmark already exists, skipping build")
		return nil
	}

	fmt.Println("Building uWebSockets HelloWorldBenchmark server...")

	// Build using GNU Make (Linux)
	makeCmd := exec.Command("make", "examples")
	makeCmd.Dir = b.uWebSocketsPath
	makeCmd.Env = append(os.Environ(), 
		"WITH_OPENSSL=0",
		"WITH_ZLIB=0",
		"WITH_LTO=0",
	)
	makeCmd.Stdout = os.Stdout
	makeCmd.Stderr = os.Stderr
	
	if err := makeCmd.Run(); err != nil {
		return fmt.Errorf("failed to build uWebSockets examples: %w", err)
	}

	// Copy the built binary to bin directory
	srcBinary := filepath.Join(b.uWebSocketsPath, "HelloWorldBenchmark")

	if err := copyFile(srcBinary, b.uWebSocketsBin); err != nil {
		return fmt.Errorf("failed to copy HelloWorldBenchmark: %w", err)
	}

	fmt.Printf("Successfully built: %s\n", b.uWebSocketsBin)
	return nil
}

func (b *Builder) GetBinaryPath() string {
	return b.outputBinary
}

func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0755)
}
