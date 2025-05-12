package network

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"tg-bot/internal/lib/api"
	"tg-bot/internal/lib/formatter"
	"tg-bot/internal/tgbot/user"
	"time"
)

var mu sync.Mutex

type Networker interface {
	AddUserRequest(chatId int64, username string) string
	AddLinkRequest(chatId int64, link, category string) string
	RemoveLinkRequest(chatId int64, link string) string
	GetUsersLinkRequest(chatId int64) string
	CheckUserRequest(chatId int64) bool
}

type Network struct {
	http      *HttpClient
	serverUrl string
}

func New(url string) *Network {
	client := &HttpClient{http.Client{Timeout: time.Second * 10}}
	return &Network{http: client, serverUrl: url + "user/"}
}

func (network *Network) AddUserRequest(chatId int64, username string) string {
	newUser := api.User{ChatID: chatId, Username: username}
	code := network.http.makePostRequest(network.serverUrl+"addUser", newUser)
	switch code {
	case http.StatusOK:
		return "Пользователь успешно добавлен"
	case http.StatusConflict:
		return "Пользователь уже существует"
	default:
		return "Произошла ошибка, попробуйте еще раз"
	}
}

func (network *Network) AddLinkRequest(chatId int64, link, category string) string {
	newUserLink := api.LinkRequest{ChatID: chatId, Link: link, Category: category}
	code := network.http.makePostRequest(network.serverUrl+"addLink", newUserLink)
	switch code {
	case http.StatusOK:
		return "Ccылка успешно добавлена!"
	case http.StatusConflict:
		return "Эта ссылка уже у вас существует!"
	default:
		return "Произошла ошибка, попробуйте еще раз"
	}
}

func (network *Network) RemoveLinkRequest(chatId int64, link string) string {
	removeLink := fmt.Sprintf("%sremoveLink?chatId=%d&link=%s", network.serverUrl, chatId, link)
	code := network.http.makeRemoveRequest(removeLink)
	switch code {
	case http.StatusOK:
		return "Ccылка успешно удалена!"
	case http.StatusNotFound:
		return "У пользователя нет такой ссылки"
	default:
		return "Произошла ошибка, попробуйте еще раз"
	}
}

func (network *Network) GetUsersLinkRequest(chatId int64) string {
	getLink := fmt.Sprintf("%sgetLinks?chatId=%d", network.serverUrl, chatId)
	code, links := network.http.makeGetRequest(getLink)

	switch code {
	case http.StatusOK:
		return formatter.FormatLinksByCategory(links)
	case http.StatusNotFound:
		return "У вас нет ссылок"
	default:
		log.Println()
		return "Произошла ошибка при получении ссылок, попробуйте позже."
	}
}

func (network *Network) CheckUserRequest(chatId int64) bool {
	userForCheck := api.User{ChatID: chatId}

	code := network.http.makePostRequest(network.serverUrl+"isExist", userForCheck)
	if code == 200 {
		if _, exist := user.Users[chatId]; !exist {
			user.Users[chatId] = &user.User{ChatId: chatId, State: user.NONE}
		}
		return true
	}
	return false
}
