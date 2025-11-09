package config

import (
	"fmt"
	"os"
)

type Config struct {
	Database DatabaseConfig
	Bot      BotConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type BotConfig struct {
	Token string
}

func (db *DatabaseConfig) GetConnectionString() string {
	connStr := os.Getenv("DB_CONN_STRING")
	if connStr != "" {
		return connStr
	}
	
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		db.Host, db.Port, db.User, db.Password, db.DBName, db.SSLMode)
}

func LoadConfig() *Config {
	connStr := os.Getenv("DB_CONN_STRING")
	token := requireEnv("BOT_TOKEN")
	
	if connStr != "" {
		return &Config{
			Database: DatabaseConfig{},
			Bot: BotConfig{
				Token: token,
			},
		}
	}
	
	return &Config{
		Database: DatabaseConfig{
			Host:     requireEnv("DB_HOST"),
			Port:     requireEnv("DB_PORT"),
			User:     requireEnv("DB_USER"),
			Password: requireEnv("DB_PASSWORD"),
			DBName:   requireEnv("DB_NAME"),
			SSLMode:  requireEnv("DB_SSLMODE"),
		},
		Bot: BotConfig{
			Token: token,
		},
	}
}

func requireEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("Required environment variable %s is not set", key))
	}
	return value
}

