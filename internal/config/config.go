package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	WhatsAppApiKey string
	Port           string
	VerifyToken    string
	PhoneID        string
	RabbitmqUrl    string
}

func Load() *Config {
	return &Config{
		WhatsAppApiKey: getEnv("WHATSAPP_ACCESS_TOKEN", "your_whatsapp_access_token"),
		Port:           getEnv("PORT", "3000"),
		VerifyToken:    getEnv("VERIFY_TOKEN", "your_verification_token"),
		PhoneID:        getEnv("PHONE_ID", "your_phone_ids"),
		RabbitmqUrl:    getEnv("RABBITMQ_URL", "your_rabbitmq_url"),
	}
}

func getEnv(key, defaultValue string) string {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	value := os.Getenv(key)
	if value != "" {
		return value
	}
	return defaultValue
}
