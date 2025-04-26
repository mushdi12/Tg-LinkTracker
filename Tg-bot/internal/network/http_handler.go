package network

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"tg-bot/internal/user"
	"time"
)

const (
	POST   = "POST"
	GET    = "GET"
	DELETE = "DELETE"
)

type HttpCode int

var (
	mu         sync.Mutex
	httpClient = http.Client{Timeout: 30 * time.Second}
	ServerURL  = "-_-"
)

func AddUserRequest(chatId int64, username string) (HttpCode, error) {
	newUser := User{ChatID: chatId, Username: username}
	code, err := makePostRequest(ServerURL+"addUser", POST, newUser)
	if err != nil {
		return 0, err
	}
	return code, nil
}

func AddLinkRequest(chatId int64, link, category string) (HttpCode, error) {
	newUserLink := LinkRequest{ChatID: chatId, Link: link, Category: category}
	code, err := makePostRequest(ServerURL+"addLink", POST, newUserLink)
	if err != nil {
		return 0, err
	}
	return code, nil
}

func CheckUser(chatId int64) bool {
	mu.Lock()
	_, exist := user.Users[chatId]
	mu.Unlock()

	if exist {
		return true
	}

	checkUser := User{ChatID: chatId}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	answer := make(chan int, 1) // буфер на 1, чтобы избежать утечки горутины

	go func() {
		code, err := makePostRequest(ServerURL+"isExist", POST, checkUser)
		if err != nil {
			answer <- -1
			return
		}
		answer <- int(code)
	}()

	select {
	case <-ctx.Done():
		return false
	case code := <-answer:
		if code != 200 {
			return false
		}

		// только если код 200 — добавляем пользователя
		mu.Lock()
		user.Users[chatId] = &user.User{ChatId: chatId, State: user.NONE}
		mu.Unlock()

		return true
	}
}

func RemoveLinkRequest(chatId int64, link string) (HttpCode, error) {
	url := fmt.Sprintf("%sremoveLink?chatId=%d&link=%s", ServerURL, chatId, link)

	req, err := http.NewRequest(DELETE, url, nil)
	if err != nil {
		log.Printf("RemoveLinkRequest err: %v\n", err)
		return 0, err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		log.Printf("RemoveLinkRequest err: %v\n", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Ошибка чтения ответа: %v", err)
		return 0, err
	}
	log.Printf("Ответ от сервера: %s", string(body))

	return HttpCode(resp.StatusCode), nil
}

func GetUsersLinkRequest(chatId int64) ([]Link, HttpCode, error) {
	fmt.Println("1234")
	url := fmt.Sprintf("%sgetLinks?chatId=%d", ServerURL, chatId)

	resp, err := httpClient.Get(url)
	if err != nil {
		log.Printf("GetUsersLinkRequest err: %v\n", err)
		return nil, 0, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Ошибка чтения ответа: %v", err)
		return nil, 0, err
	}

	log.Printf("Ответ от сервера: %s", string(body))

	var response LinksResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("Ошибка при разборе ответа: %v", err)
		return nil, 0, err
	}
	fmt.Println("12345")
	return response.Links, HttpCode(resp.StatusCode), nil
}

func makePostRequest(url string, httpMethod string, dto any) (HttpCode, error) {
	var req *http.Request
	var err error

	data, err := json.Marshal(dto)
	if err != nil {
		log.Printf("Ошибка сериализации структуры: %v", err)
		return 0, err
	}

	req, err = http.NewRequest(httpMethod, url, bytes.NewBuffer(data))
	if err != nil {
		log.Printf("Ошибка создания запроса: %v", err)
		return 0, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		log.Printf("Ошибка отправки запроса: %v", err)
		return 0, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Ошибка чтения ответа: %v", err)
		return 0, err
	}
	log.Printf("Ответ от сервера: %s", string(body))

	return HttpCode(resp.StatusCode), nil
}
