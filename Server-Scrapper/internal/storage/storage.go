package storage

type StorageContext struct {
	ChatID   int64
	UserName string
	Link     UserLink
}

type UserLinks []UserLink

type UserLink struct {
	LinkUrl      string
	LinkCategory string
}
