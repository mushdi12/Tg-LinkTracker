package repositories

import "time"

type WatchedRepo struct {
	ID        int64
	ChatID    int64
	RepoOwner string
	RepoName  string
	LastSHA   string
	LastIssue int
	LastPR    int
	UpdatedAt *time.Time
}

type RepoInfo struct {
	RepoOwner string
	RepoName  string
	LastSHA   string
	LastIssue int
	LastPR    int
}
