package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

const (
	T_API string = "https://api.telegram.org/bot"
	W_API string = "https://api.openweathermap.org/data/2.5"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	BOT_TOKEN := os.Getenv("BOT_TOKEN")
	OWM_TOKEN := os.Getenv("OWM_TOKEN")
	URL := T_API + BOT_TOKEN
	OFFSET := 0

	for {
		updates, err := GetUpdates(URL, OFFSET)
		if err != nil {
			log.Fatal(err)
		}

		for _, update := range updates {
			log.Printf("[%s]: %s", update.Message.From.Username, update.Message.Text)

			err := Reply(URL, update, OWM_TOKEN)
			if err != nil {
				log.Fatal(err)
			}

			OFFSET = update.UpdateID + 1
		}
	}
}

func GetUpdates(URL string, OFFSET int) (updates []Update, err error) {
	resp, err := http.Get(URL + "/getUpdates" + "?offset=" + strconv.Itoa(OFFSET))
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	response := new(Response)
	json.Unmarshal(body, &response)

	updates = response.Result

	return
}

func GetReplyData(resp *http.Response) (replyData string, err error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	data := new(Data)
	json.Unmarshal(body, &data)

	if fmt.Sprintf("%v", data.Cod) != "200" {
		replyData = data.Message
	} else {
		replyData = fmt.Sprintf("City: %v\nWeather: %v\nTemperature: %v°C\nWind: %.1fkm/h\n", data.Name, data.Weather[0].Main, math.Round(data.Main.Temp-273.15), data.Wind.Speed*2.449)
	}

	return
}

func Reply(URL string, update Update, OWM_TOKEN string) (err error) {
	replyMessage := new(ReplyMessage)
	replyMessage.ChatID = update.Message.Chat.ID

	switch update.Message.Text {
	case "/start", "/help":
		replyMessage.Text = fmt.Sprintln("This is a Weather Bot\n\nEnter the city name to see the weather forecast")
	default:
		resp, r_err := http.Get(W_API + "/weather?q=" + update.Message.Text + "&appid=" + OWM_TOKEN)
		if r_err != nil {
			return
		}
		defer resp.Body.Close()

		replyMessage.Text, err = GetReplyData(resp)
		if err != nil {
			return
		}
	}

	buffer, _ := json.Marshal(replyMessage)

	_, err = http.Post(URL+"/sendMessage", "application/json", bytes.NewBuffer(buffer))
	if err != nil {
		return
	}

	return
}
