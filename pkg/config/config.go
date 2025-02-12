package config

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"

	"forum/pkg/logger"
)

// Config holds application configuration
type Config struct {
	Port int
}

// LoadConfig reads configuration from a .env file
func LoadConfig() *Config {
	// Load environment variables from file
	err := loadEnvFile("./pkg/config/variables.env")
	if err != nil {
		logger.LogWithDetails(err)
		log.Fatal(err)
	}

	port, err := strconv.Atoi(getEnv("PORT", "8080"))
	if err != nil {
		logger.LogWithDetails(err)
		log.Fatalf("Invalid PORT value: %v", err)
	}

	return &Config{
		Port: port,
	}
}

func loadEnvFile(filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		logger.LogWithDetails(err)
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		os.Setenv(key, value)
	}

	if err := scanner.Err(); err != nil {
		logger.LogWithDetails(err)
		return err
	}
	return nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
