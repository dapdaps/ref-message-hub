package model

type MessageParam struct {
	TgGroup string `json:"tg_group"`
	Slack   string `json:"slack"`
	Email   string `json:"email"`
	Content string `json:"content" binding:"required"`
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
