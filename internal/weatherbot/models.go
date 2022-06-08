package weatherbot

type Response struct {
	Result []Update `json:"result,omitempty"`
}

type Update struct {
	UpdateID int     `json:"update_id,omitempty"`
	Message  Message `json:"message,omitempty"`
}

type Message struct {
	Chat Chat   `json:"chat,omitempty"`
	From User   `json:"from,omitempty"`
	Text string `json:"text,omitempty"`
}

type Chat struct {
	ID int `json:"id,omitempty"`
}

type User struct {
	Username string `json:"username,omitempty"`
}

type ReplyMessage struct {
	ChatID int    `json:"chat_id,omitempty"`
	Text   string `json:"text,omitempty"`
}

type Data struct {
	Cod     interface{} `json:"cod,omitempty"`
	Message string      `json:"message,omitempty"`
	Name    string      `json:"name,omitempty"`
	Weather []Weather   `json:"weather,omitempty"`
	Main    Main        `json:"main,omitempty"`
	Wind    Wind        `json:"wind,omitempty"`
}

type Weather struct {
	Main string `json:"main,omitempty"`
}

type Main struct {
	Temp float64 `json:"temp,omitempty"`
}

type Wind struct {
	Speed float64 `json:"speed,omitempty"`
}
