package network

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"tg-bot/internal/lib/api"
)

const (
	POST   = "POST"
	GET    = "GET"
	DELETE = "DELETE"
)

type HttpCode int

type HttpInterface interface {
	makePostRequest(url string, dto any) HttpCode
	makeRemoveRequest(url string) HttpCode
	makeGetRequest(url string) (HttpCode, []api.Link)
}

type HttpClient struct {
	http.Client
}

func (httpClient *HttpClient) makePostRequest(url string, dto any) HttpCode {
	op := "lib.api.makePostRequest"
	var req *http.Request

	data, err := json.Marshal(dto)
	if err != nil {
		_ = fmt.Errorf("%s: %w", op, err)
		return 0
	}

	req, err = http.NewRequest(POST, url, bytes.NewBuffer(data))
	if err != nil {
		_ = fmt.Errorf("%s: %w", op, err)
		return 0
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		_ = fmt.Errorf("%s: %w", op, err)
		return 0
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		_ = fmt.Errorf("%s: %w", op, err)
		return 0
	}

	log.Printf("Ответ от сервера: %s", string(body))

	return HttpCode(resp.StatusCode)
}

func (httpClient *HttpClient) makeRemoveRequest(url string) HttpCode {
	op := "lib.api.makeRemoveRequest"
	var req *http.Request
	req, err := http.NewRequest(DELETE, url, nil)
	if err != nil {
		_ = fmt.Errorf("%s: %w", op, err)
		return 0
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		_ = fmt.Errorf("%s: %w", op, err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		_ = fmt.Errorf("%s: %w", op, err)
		return 0
	}
	log.Printf("Ответ от сервера: %s", string(body))

	return HttpCode(resp.StatusCode)
}

func (httpClient *HttpClient) makeGetRequest(url string) (HttpCode, []api.Link) {
	op := "lib.api.makeGetRequest"

	resp, err := httpClient.Get(url)
	if err != nil {
		_ = fmt.Errorf("%s: %w", op, err)
		return 0, nil
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, nil
	}

	log.Printf("Ответ от сервера: %s", string(body))

	var response api.LinksResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return 0, nil
	}
	return HttpCode(resp.StatusCode), response.Links
}
