package telegram

import (
	"net/http"
	"server-scrapper/internal/config"
	"server-scrapper/internal/lib/api"
	"server-scrapper/internal/lib/network"
	"server-scrapper/internal/lib/repositories"
	"time"
)

type Networker interface {
	SendRepoUpdates(repo repositories.WatchedRepo, message string)
}

type Telegram struct {
	http      *network.HttpClient
	serverUrl string
}

func New(cfg config.TelegramBot) *Telegram {
	client := &network.HttpClient{Client: http.Client{Timeout: time.Second * 10}}
	return &Telegram{http: client, serverUrl: cfg.URL}
}

func (net *Telegram) SendRepoUpdates(repo repositories.WatchedRepo, message string) {
	data := api.GitRepo{To: repo.ChatID, Name: repo.RepoName, Owner: repo.RepoOwner, Msg: message}
	net.http.MakePostRequest(net.serverUrl+"githubInfo", data)
}
