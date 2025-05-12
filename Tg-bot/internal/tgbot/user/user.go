package user

import (
	"errors"
	"sync"
)

var (
	Users = make(map[int64]*User)
	mu    sync.Mutex
)

type UserInterface interface {
	getState() int
}

type User struct {
	ChatId   int64
	State    int
	Link     string
	Filter   string
	Category string
}

func GetState(chatId int64) (int, error) {
	mu.Lock()
	defer mu.Unlock()
	user, exist := Users[chatId]
	if !exist {
		return 0, errors.New("Not exist")
	}
	return user.State, nil
}

func ResetState(chatId int64) {
	mu.Lock()
	defer mu.Unlock()
	user, _ := Users[chatId]
	user.State = NONE
}
