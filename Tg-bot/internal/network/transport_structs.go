package network

type User struct {
	ChatID   int64  `json:"chat_id" gorm:"column:chat_id"`
	Username string `json:"username" gorm:"column:username"`
}

type LinkRequest struct {
	ChatId   int64  `json:"chat_id"`
	Link     string `json:"link" gorm:"column:irl"`
	Category string `json:"category"`
	Filters  string `json:"filters"`
}

type LinksResponse struct {
	Link []string `json:"links"`
}
