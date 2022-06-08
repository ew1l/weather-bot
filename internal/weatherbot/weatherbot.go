package weatherbot

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
)

const (
	TelegramAPI       = "https://api.telegram.org/bot"
	OpenWeatherMapAPI = "https://api.openweathermap.org/data/2.5"
)

var OpenWeatherMapToken string

type WeatherBot struct {
	URL    string
	Offset int
}

func New() *WeatherBot {
	OpenWeatherMapToken = os.Getenv("OWM_TOKEN")
	return &WeatherBot{
		URL:    TelegramAPI + os.Getenv("BOT_TOKEN"),
		Offset: 0,
	}
}

func (wb *WeatherBot) Start() {
	for {
		updates, err := wb.GetUpdates()
		if err != nil {
			log.Println(err)
		}

		for _, update := range updates {
			log.Printf("[%s]: %s", update.Message.From.Username, update.Message.Text)

			if err := wb.Reply(update); err != nil {
				log.Println(err)
			}

			wb.Offset = update.UpdateID + 1
		}
	}
}

func (wb *WeatherBot) GetUpdates() ([]Update, error) {
	resp, err := http.Get(wb.URL + "/getUpdates" + "?offset=" + strconv.Itoa(wb.Offset))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	response := new(Response)
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	updates := response.Result

	return updates, nil
}

func (wb *WeatherBot) GetReplyData(resp *http.Response) (string, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	data := new(Data)
	if err := json.Unmarshal(body, &data); err != nil {
		return "", err
	}

	var replyData string
	switch fmt.Sprintf("%v", data.Cod) {
	case "200":
		replyData = fmt.Sprintf("City: %v\nWeather: %v\nTemperature: %v°C\nWind: %.1fkm/h\n",
			data.Name,
			data.Weather[0].Main,
			math.Round(data.Main.Temp-273.15),
			data.Wind.Speed*2.449,
		)
	default:
		replyData = data.Message
	}

	return replyData, nil
}

func (wb *WeatherBot) Reply(update Update) error {
	replyMessage := new(ReplyMessage)
	replyMessage.ChatID = update.Message.Chat.ID

	switch update.Message.Text {
	case "/start", "/help":
		replyMessage.Text = fmt.Sprintln("Weather Bot\n\nEnter the city name to see the weather forecast")
	default:
		resp, err := http.Get(OpenWeatherMapAPI + "/weather?q=" + update.Message.Text + "&appid=" + OpenWeatherMapToken)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		replyMessage.Text, err = wb.GetReplyData(resp)
		if err != nil {
			return err
		}
	}

	buffer, err := json.Marshal(replyMessage)
	if err != nil {
		return err
	}

	if _, err := http.Post(wb.URL+"/sendMessage", "application/json", bytes.NewBuffer(buffer)); err != nil {
		return err
	}

	return nil
}
