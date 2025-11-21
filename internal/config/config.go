package config

import "path/filepath"

type ServerConfig struct {
	Name     string
	Dir      string
	StartCmd []string
	Port     int
}

func GetServers(baseDir string) []ServerConfig {
	binDir := filepath.Join(baseDir, "bin")
	uWebSocketsBin := filepath.Join(binDir, "HelloWorldBenchmark")

	return []ServerConfig{
		{
			Name:     "go-http",
			Dir:      filepath.Join(baseDir, "api", "go-http"),
			StartCmd: []string{"go", "run", "main.go"},
			Port:     8080,
		},
		{
			Name:     "go-fasthttp",
			Dir:      filepath.Join(baseDir, "api", "go-fasthttp"),
			StartCmd: []string{"go", "run", "main.go"},
			Port:     8080,
		},
		{
			Name:     "python-fastapi",
			Dir:      filepath.Join(baseDir, "api", "python-fastapi"),
			StartCmd: []string{"python", "main.py"},
			Port:     8080,
		},
		{
			Name:     "node-http",
			Dir:      filepath.Join(baseDir, "api", "node-http"),
			StartCmd: []string{"node", "index.js"},
			Port:     8080,
		},
		{
			Name:     "nginx-static",
			Dir:      filepath.Join(baseDir, "api", "nginx-static"),
			StartCmd: []string{"nginx", "-p", ".", "-c", "nginx.conf"},
			Port:     8080,
		},
		{
			Name:     "cpp-uwebsockets",
			Dir:      binDir,
			StartCmd: []string{uWebSocketsBin},
			Port:     8080,
		},
	}
}

func GetServerByName(baseDir, name string) *ServerConfig {
	servers := GetServers(baseDir)
	for _, s := range servers {
		if s.Name == name {
			return &s
		}
	}
	return nil
}
