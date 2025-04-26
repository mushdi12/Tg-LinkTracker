package api

type LinkRequest struct {
	ChatID   int64  `json:"chat_id"`
	Link     string `json:"link" `
	Category string `json:"category"`
}

type User struct {
	ChatID   int64  `json:"chat_id" `
	Username string `json:"username" `
}
