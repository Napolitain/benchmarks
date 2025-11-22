package server

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"time"

	"github.com/benchmarks/internal/config"
)

type Server struct {
	config *config.ServerConfig
	cmd    *exec.Cmd
}

func New(cfg *config.ServerConfig) *Server {
	return &Server{config: cfg}
}

func (s *Server) GetPID() int {
	if s.cmd != nil && s.cmd.Process != nil {
		return s.cmd.Process.Pid
	}
	return -1
}

func (s *Server) Start() error {
	fmt.Printf("Starting %s...\n", s.config.Name)

	s.cmd = exec.Command(s.config.StartCmd[0], s.config.StartCmd[1:]...)
	s.cmd.Dir = s.config.Dir
	s.cmd.Stdout = os.Stdout
	s.cmd.Stderr = os.Stderr

	if err := s.cmd.Start(); err != nil {
		return fmt.Errorf("failed to start %s: %w", s.config.Name, err)
	}

	// Wait for server to be ready
	fmt.Printf("Waiting for %s to be ready...\n", s.config.Name)
	if err := s.waitForServer(5 * time.Second); err != nil {
		s.Stop()
		return fmt.Errorf("%s failed to start: %w", s.config.Name, err)
	}

	fmt.Printf("%s is ready\n", s.config.Name)
	return nil
}

func (s *Server) Stop() error {
	if s.cmd == nil || s.cmd.Process == nil {
		return nil
	}

	fmt.Printf("Stopping %s...\n", s.config.Name)

	// Special handling for nginx
	if s.config.Name == "nginx-static" {
		stopCmd := exec.Command("nginx", "-p", ".", "-c", "nginx.conf", "-s", "stop")
		stopCmd.Dir = s.config.Dir
		stopCmd.Run()
		time.Sleep(1 * time.Second)
		return nil
	}

	// Kill process (Linux)
	s.cmd.Process.Kill()
	s.cmd.Wait()
	time.Sleep(2 * time.Second)
	return nil
}

func (s *Server) waitForServer(timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	addr := fmt.Sprintf("localhost:%d", s.config.Port)

	for time.Now().Before(deadline) {
		conn, err := net.DialTimeout("tcp", addr, 500*time.Millisecond)
		if err == nil {
			conn.Close()
			return nil
		}
		time.Sleep(200 * time.Millisecond)
	}

	return fmt.Errorf("timeout waiting for server on port %d", s.config.Port)
}
