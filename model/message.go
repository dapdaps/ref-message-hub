package model

type MessageParam struct {
	Title    string `json:"title"`
	Type     string `json:"type" binding:"required"`
	Product  string `json:"product" binding:"required"`
	Content  string `json:"content" binding:"required"`
	Level    string `json:"level" binding:"required"`
	Telegram bool   `json:"telegram"`
	Slack    bool   `json:"slack"`
	Email    bool   `json:"email"`
}

type MessageResponse struct {
	Telegram bool `json:"telegram"`
	Slack    bool `json:"slack"`
	Email    bool `json:"email"`
}

type SlackMessageBlock struct {
	Type string       `json:"type"`
	Text SlackMessage `json:"text"`
}

type SlackMessage struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type TelegramMessage struct {
	ChatId int64  `json:"chat_id"`
	Text   string `json:"text"`
}

type TelegramResponse struct {
	Ok bool `json:"ok"`
}
