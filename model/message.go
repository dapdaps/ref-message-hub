package model

type MessageParam struct {
	Title   string `json:"title"`
	Content string `json:"content" binding:"required"`
	TgGroup string `json:"tg_group"`
	Slack   string `json:"slack"`
	Email   bool   `json:"email"`
	Level   string `json:"level"`
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
