package network

type User struct {
	ChatID   int64  `json:"chat_id" `
	Username string `json:"username" `
}

type LinkRequest struct {
	ChatID   int64  `json:"chat_id"`
	Link     string `json:"link" gorm:"column:irl"`
	Category string `json:"category"`
}

type LinksResponse struct {
	Links []Link `json:"links"`
}

type Link struct {
	URL      string `json:"irl"`
	Category string `json:"category"`
}
