package env

import (
	"ft-bot/store"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
	"strings"
)

var E *store.Environment

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
	E = NewEnvironment()
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

// Simple helper function to read an environment variable into integer or return a default value
func getEnvAsInt(name string, defaultVal int) int {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultVal
}

// Helper to read an environment variable into a bool or return default value
func getEnvAsBool(name string, defaultVal bool) bool {
	valStr := getEnv(name, "")
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}
	return defaultVal
}

// Helper to read an environment variable into a string slice or return default value
func getEnvAsSlice(name string, defaultVal []string, sep string) []string {
	valStr := getEnv(name, "")

	if valStr == "" {
		return defaultVal
	}

	val := strings.Split(valStr, sep)

	return val
}

func NewEnvironment() *store.Environment {
	return &store.Environment{
		Token:         getEnv("TOKEN", ""),
		BotPrefix:     getEnv("BOTPREFIX", "!"),
		MySqlHost:     getEnv("MYSQL_HOST", ""),
		MySqlPort:     getEnvAsInt("MYSQL_PORT", 3306),
		MySqlDatabase: getEnv("MYSQL_DATABASE", ""),
		MySqlUser:     getEnv("MYSQL_USER", ""),
		MySqlPassword: getEnv("MYSQL_PASSWORD", ""),
		GuildId:       getEnv("GUILD_ID", ""),
		AdminRoles:    getEnvAsSlice("ADMIN_ROLES", []string{}, ","),
	}
}
