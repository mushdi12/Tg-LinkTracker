package network

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	POST = "POST"
	GET  = "GET"
)

type HttpCode int

type HttpWorker interface {
	MakeGetRequest(url string) ([]byte, error)
	MakePostRequest(url string, dto any) HttpCode
}

type HttpClient struct {
	http.Client
}

func (httpClient *HttpClient) MakeGetRequest(url string) ([]byte, error) {
	req, err := http.NewRequest(GET, url, nil)
	if err != nil {
		return nil, err
	}

	// GitHub требует заголовок User-Agent
	req.Header.Set("User-Agent", "go-github-watcher")

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET %s: статус %d", url, resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

func (httpClient *HttpClient) MakePostRequest(url string, dto any) HttpCode {
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
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		_ = fmt.Errorf("%s: %w", op, err)
		return 0
	}

	log.Printf("Ответ от сервера: %s", string(body))

	return HttpCode(resp.StatusCode)
}
