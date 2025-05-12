package utils

import (
	"sync"
	"tg-bot/internal/tgbot/user"
)

var mu sync.Mutex

// возвращает false если во время команды использутся другая
func CheckDoubleCmmd(chatId int64) bool {
	mu.Lock()
	defer mu.Unlock()

	u, _ := user.Users[chatId]
	if u == nil {
		return true
	}

	if u.State > 1 {
		return false
	}
	return true
}
