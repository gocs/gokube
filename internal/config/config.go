package config

import (
	"log"
	"os"
)

var (
	YoutubeAPIKey string
)

func init() {
	YoutubeAPIKey = os.Getenv("YOUTUBE_API_KEY")
	if YoutubeAPIKey == "" {
		log.Fatal("YOUTUBE_API_KEY is not set")
	}
}
