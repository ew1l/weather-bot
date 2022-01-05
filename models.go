package main

type Response struct {
	Result []Update `json:"result"`
}

type Update struct {
	UpdateID int     `json:"update_id"`
	Message  Message `json:"message"`
}

type Message struct {
	Chat Chat   `json:"chat"`
	From User   `json:"from"`
	Text string `json:"text"`
}

type Chat struct {
	ID int `json:"id"`
}

type User struct {
	Username string `json:"username"`
}

type ReplyMessage struct {
	ChatID int    `json:"chat_id"`
	Text   string `json:"text"`
}

type Data struct {
	Cod     interface{} `json:"cod"`
	Message string      `json:"message"`
	Name    string      `json:"name"`
	Weather []Weather   `json:"weather"`
	Main    Main        `json:"main"`
	Wind    Wind        `json:"wind"`
}

type Weather struct {
	Main string `json:"main"`
}

type Main struct {
	Temp float64 `json:"temp"`
}

type Wind struct {
	Speed float64 `json:"speed"`
}
