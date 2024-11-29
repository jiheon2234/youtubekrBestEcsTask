package config

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	VideoApiURL   string
	CommentApiURL string
	ApiKey        string
	PgDsn         string
	GoRoutineCnt  int
}

func NewConfig() *Config {

	apikey := mustGetEnv("YOUTUBE_API_KEY")
	//log.Fatalf(apikey)
	pgdsn := mustGetEnv("PGDSN")
	goRoutineCnt, err := strconv.Atoi(mustGetEnv("GOROUTINE_CNT"))
	if err != nil {
		panic(err)
	}

	return &Config{
		VideoApiURL:   "https://www.googleapis.com/youtube/v3/videos",
		CommentApiURL: "https://youtube.googleapis.com/youtube/v3/commentThreads",
		ApiKey:        apikey,
		PgDsn:         pgdsn,
		GoRoutineCnt:  goRoutineCnt,
	}
}

func mustGetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Environment variable %s is required but not set", key)
	}
	return value
}
