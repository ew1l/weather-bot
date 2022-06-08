package main

import (
	"log"

	"github.com/ew1l/weather-bot/internal/weatherbot"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	bot := weatherbot.New()
	bot.Start()
}
