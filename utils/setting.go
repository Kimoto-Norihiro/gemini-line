package utils

import (
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

type Setting struct {
	Port               string
	DatabaseUrl        string
	ChannelSecret      string
	ChannelAccessToken string
	GeminiApiKey       string
}

func LoadSetting() (*Setting, error) {
	path := filepath.Join(".env")
	err := godotenv.Load(path)
	if err != nil {
		return nil, err
	}

	return &Setting{
		Port:               os.Getenv("PORT"),
		DatabaseUrl:        os.Getenv("DATABASE_URL"),
		ChannelSecret:      os.Getenv("CHANNEL_SECRET"),
		ChannelAccessToken: os.Getenv("CHANNEL_ACCESS_TOKEN"),
		GeminiApiKey:       os.Getenv("GEMINI_API_KEY"),
	}, nil
}
