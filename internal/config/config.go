package config

import "os"

type Config struct {
	Port       string
	BackendURL string
}

// Load all env configs
func Load() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	backendUrl := os.Getenv("BACKEND_URL")
	if backendUrl == "" {
		backendUrl = "localhost:8081"
	}

	return &Config{
		Port:       port,
		BackendURL: backendUrl,
	}
}
