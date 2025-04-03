package config

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort    string
	DBPath        string
	BotToken      string
	BaseURL       string
	PosterToken   string
	EncryptionKey []byte
	Admins        []int
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load()

	config := &Config{
		ServerPort:    getEnv("PORT", "8080"),
		DBPath:        getEnv("DB_PATH", "registration.db?mode=rwc"),
		BotToken:      getEnv("BOT_TOKEN", ""),
		BaseURL:       getEnv("BASE_URL", ""),
		PosterToken:   getEnv("POSTER_TOKEN", ""),
		EncryptionKey: []byte(getEnv("ENCRYPTION_KEY", "")),
	}

	adminsStr := getEnv("ADMINS", "")
	config.Admins = parseAdmins(adminsStr)

	// Проверка обязательных переменных
	if config.BotToken == "" {
		return nil, fmt.Errorf("BOT_TOKEN не задан")
	}
	if config.BaseURL == "" {
		return nil, fmt.Errorf("BASE_URL не задан")
	}
	if config.PosterToken == "" {
		return nil, fmt.Errorf("POSTER_TOKEN не задан")
	}
	if len(config.EncryptionKey) == 0 {
		return nil, fmt.Errorf("ENCRYPTION_KEY не задан")
	}

	return config, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func parseAdmins(adminsStr string) []int {
	var admins []int
	if adminsStr == "" {
		return admins
	}

	parts := strings.Split(adminsStr, ",")
	for _, p := range parts {
		admins = append(admins, parseInt(p))
	}
	return admins
}

func parseInt(s string) int {
	var i int
	_, err := fmt.Sscanf(s, "%d", &i)
	if err != nil {
		log.Printf("Ошибка парсинга int: %v", err)
	}
	return i
}
