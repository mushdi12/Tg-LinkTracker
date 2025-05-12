package watcher

import (
	"context"
	"fmt"
	"github.com/go-co-op/gocron"
	"log"
	"server-scrapper/internal/config"
	"server-scrapper/internal/lib/github"
	"server-scrapper/internal/lib/repositories"
	"server-scrapper/internal/lib/telegram"
	"server-scrapper/internal/storage/postgres"
	"time"
)

type Watcher struct {
	*gocron.Scheduler
	storage   *postgres.Storage
	tgClient  *telegram.Telegram
	gitClient *github.GitHub
	interval  int
	count     int
}

func New(storage *postgres.Storage, tgClient *telegram.Telegram, gitClient *github.GitHub, cfg config.Watcher) *Watcher {
	watcher := &Watcher{
		Scheduler: gocron.NewScheduler(time.UTC),
		storage:   storage,
		tgClient:  tgClient,
		gitClient: gitClient,
		interval:  cfg.IntervalMinutes,
		count:     cfg.MaxReposPerCheck,
	}
	return watcher
}

func (watcher *Watcher) StartMonitoring() {
	op := "watcher.watcher.StartMonitoring"
	_, err := watcher.Every(watcher.interval).Minute().Do(func() {
		monitoringRepositories(watcher)
	})

	if err != nil {
		log.Printf("%s: %w", op, err)
	}

	watcher.StartAsync()
}

// функция горутина
func monitoringRepositories(watcher *Watcher) {
	op := "watcher.watcher.monitoringRepositories"

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	repos, err := watcher.storage.GetOldestRepos(ctx, watcher.count)
	if err != nil {
		fmt.Printf("%s: get repos: %s\n", op, err)
		return
	}

	for _, repo := range repos {
		go handleRepository(repo, watcher)
	}
}

func handleRepository(repo repositories.WatchedRepo, watcher *Watcher) {
	op := "watcher.watcher.handleRepository"
	isUpdate, message := watcher.gitClient.CheckGitData(&repo)

	if isUpdate {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := watcher.storage.UpdateRepo(ctx, repo)
		if err != nil {
			fmt.Printf("%s: %v\n", op, err)
		}

		watcher.tgClient.SendRepoUpdates(repo, message)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := watcher.storage.UpdateRepoCheckTime(ctx, repo.ID)
	if err != nil {
		fmt.Printf("%s: %v\n", op, err)
	}
}
