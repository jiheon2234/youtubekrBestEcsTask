package config

import (
	"log"
	"os"
)

type Config struct {
	VideoApiURL   string
	CommentApiURL string
	ApiKey        string
	PgDsn         string
}

func NewConfig() *Config {

	apikey := mustGetEnv("YOUTUBE_API_KEY")
	//log.Fatalf(apikey)
	pgdsn := mustGetEnv("PGDSN")

	return &Config{
		VideoApiURL:   "https://www.googleapis.com/youtube/v3/videos",
		CommentApiURL: "https://youtube.googleapis.com/youtube/v3/commentThreads",
		ApiKey:        apikey,
		PgDsn:         pgdsn,
	}
}

func mustGetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Environment variable %s is required but not set", key)
	}
	return value
}
