package config

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

// Config holds application configuration
type Config struct {
	Port      int
	DBHost    string
	DBPort    int
	DBUser    string
	DBPass    string
	JWTSecret string
}

// LoadConfig reads configuration from a .env file
func LoadConfig() *Config {
	// Load environment variables from file
	err := loadEnvFile("./pkg/config/variables.env")
	if err != nil {
		log.Fatal(err)
	}

	port, err := strconv.Atoi(getEnv("PORT", "8080"))
	if err != nil {
		log.Fatalf("Invalid PORT value: %v", err)
	}

	dbPort, err := strconv.Atoi(getEnv("DB_PORT", "5432"))
	if err != nil {
		log.Fatalf("Invalid DB_PORT value: %v", err)
	}

	return &Config{
		Port:      port,
		DBHost:    getEnv("DB_HOST", "localhost"),
		DBPort:    dbPort,
		DBUser:    getEnv("DB_USER", "user"),
		DBPass:    getEnv("DB_PASS", "password"),
		JWTSecret: getEnv("JWT_SECRET", "mysecret"),
	}
}

func loadEnvFile(filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		// log.Fatalf("Error opening .env file: %v", err)
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
