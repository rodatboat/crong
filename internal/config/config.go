package config

import (
	"fmt"
	"os"
)

type Config struct {
	// Server
	Port string

	// Database
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
	DBSchema   string

	// Cleanup Job
	CleanupSchedule  string // Cron format, defaults to weekly (0 0 * * 0 = Sunday midnight)
	CleanupRetention int    // Days to retain execution history, defaults to 30
}

// Load reads configuration from environment variables
func Load() *Config {
	return &Config{
		Port:             getEnv("PORT", "3000"),
		DBHost:           getEnv("DB_HOST", "localhost"),
		DBPort:           getEnv("DB_PORT", "5435"),
		DBUser:           getEnv("DB_USER", "postgres"),
		DBPassword:       getEnv("DB_PASSWORD", "password"),
		DBName:           getEnv("DB_NAME", "postgres"),
		DBSSLMode:        getEnv("DB_SSL_MODE", "disable"),
		DBSchema:         getEnv("DB_SCHEMA", "crong"),
		CleanupSchedule:  getEnv("CLEANUP_SCHEDULE", "0 0 * * 0"), // Sunday midnight
		CleanupRetention: getEnvInt("CLEANUP_RETENTION", 30),      // 30 days
	}
}

// GetDSN returns the database connection string
func (c *Config) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s search_path=%s",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.DBSSLMode, c.DBSchema,
	)
}

// Helper function to get environment variable with default
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// Helper function to get environment variable as integer with default
func getEnvInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intVal, err := parseToInt(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

// Helper function to parse string to int
func parseToInt(s string) (int, error) {
	var num int
	_, err := fmt.Sscanf(s, "%d", &num)
	return num, err
}
