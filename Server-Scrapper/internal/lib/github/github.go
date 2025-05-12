package github

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"server-scrapper/internal/lib/network"
	"server-scrapper/internal/lib/repositories"
	"time"
)

type GitHub struct {
	http *network.HttpClient
}

func New() *GitHub {
	client := &network.HttpClient{Client: http.Client{Timeout: time.Second * 10}}
	return &GitHub{http: client}
}

func (gitNetwork *GitHub) CheckGitData(repo *repositories.WatchedRepo) (bool, string) {
	updated := false

	sha, commit := gitNetwork.CheckLatestCommit(repo.RepoOwner, repo.RepoName)
	if sha != repo.LastSHA {
		updated = true
		repo.LastSHA = sha
		return updated, fmt.Sprintf("📦 Новый коммит в https://github.com/%s/%s:\nСообщение: %s", repo.RepoOwner, repo.RepoName, commit)
	}

	issues := gitNetwork.CheckLatestIssues(repo.RepoOwner, repo.RepoName)
	for _, issue := range issues {
		if issue["pull_request"] != nil {
			continue
		}
		id := int(issue["number"].(float64))
		if id > repo.LastIssue {
			repo.LastIssue = id
			updated = true
			return updated, fmt.Sprintf("📝 Новый issue в https://github.com/%s/%s: #%d", repo.RepoOwner, repo.RepoName, id)
		}
		break
	}

	prs := gitNetwork.CheckLatestPR(repo.RepoOwner, repo.RepoName)
	for _, pr := range prs {
		prID := int(pr["number"].(float64))
		if prID > repo.LastPR {
			repo.LastPR = prID
			return updated, fmt.Sprintf("🚀 Новый PR в https://github.com/%s/%s: #%d", repo.RepoOwner, repo.RepoName, prID)
		}
		break
	}

	return updated, ""
}

func (gitNetwork *GitHub) CheckLatestCommit(repoOwner, repoName string) (string, string) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/commits", repoOwner, repoName)
	data, err := gitNetwork.http.MakeGetRequest(url)
	if err != nil {
		log.Println(err)
		return "", ""
	}
	sha, commit := parseCommit(data)
	return sha, commit
}

func (gitNetwork *GitHub) CheckLatestIssues(repoOwner, repoName string) []map[string]interface{} {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/issues", repoOwner, repoName)
	data, err := gitNetwork.http.MakeGetRequest(url)
	if err != nil {
		log.Println(err)
		return nil
	}
	issues := parseIssue(data)
	return issues

}

func (gitNetwork *GitHub) CheckLatestPR(repoOwner, repoName string) []map[string]interface{} {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls", repoOwner, repoName)
	data, err := gitNetwork.http.MakeGetRequest(url)
	if err != nil {
		log.Println(err)
		return nil
	}

	return parsePR(data)
}

func (gitNetwork *GitHub) GetFirstRepoInfo(parentCtx context.Context, url string, timeout time.Duration) (*repositories.RepoInfo, error) {
	const op = "network.GetFirstRepoInfo"

	ctx, cancel := context.WithTimeout(parentCtx, timeout)
	defer cancel()

	ch := make(chan *repositories.RepoInfo, 1)
	errCh := make(chan error, 1)

	go func() {
		info, err := gitNetwork.FetchRepoInfo(url)
		if err != nil {
			errCh <- err
			return
		}
		ch <- info
	}()

	select {
	case info := <-ch:
		return info, nil
	case err := <-errCh:
		return nil, fmt.Errorf("%s: %w", op, err)
	case <-ctx.Done():
		return nil, fmt.Errorf("%s: timeout", op)
	}
}

func (gitNetwork *GitHub) FetchRepoInfo(url string) (*repositories.RepoInfo, error) {
	owner, repo, err := parseRepoURL(url)
	if err != nil {
		return nil, err
	}

	info := &repositories.RepoInfo{
		RepoOwner: owner,
		RepoName:  repo,
	}

	sha, _ := gitNetwork.CheckLatestCommit(owner, repo)
	info.LastSHA = sha

	issues := gitNetwork.CheckLatestIssues(owner, repo)
	for _, issue := range issues {
		if issue["pull_request"] != nil {
			continue
		}
		id := int(issue["number"].(float64))
		if id > info.LastIssue {
			info.LastIssue = id
		}
		break
	}

	prs := gitNetwork.CheckLatestPR(owner, repo)
	for _, pr := range prs {
		prID := int(pr["number"].(float64))
		if prID > info.LastPR {
			info.LastPR = prID

		}
		break
	}

	return info, nil
}

// old version FetchRepoInfo
//func FetchRepoInfo(url string) (*repositories.RepoInfo, error) {
//	owner, repo, err := parseRepoURL(url)
//	if err != nil {
//		return nil, err
//	}
//
//	info := &repositories.RepoInfo{
//		RepoOwner: owner,
//		RepoName:  repo,
//	}
//
//	// Последний коммит
//	var commits []map[string]interface{}
//	if err := GetJSON(
//		fmt.Sprintf("https://api.github.com/repos/%s/%s/commits", owner, repo), &commits); err == nil && len(commits) > 0 {
//		if sha, ok := commits[0]["sha"].(string); ok {
//			sha, _ := CheckLatestCommit(owner, repo)
//			info.LastSHA = sha
//		}
//	}
//
//	// Последний issue
//	var issues []map[string]interface{}
//	if err := GetJSON(fmt.Sprintf("https://api.github.com/repos/%s/%s/issues", owner, repo), &issues); err == nil {
//		for _, issue := range issues {
//			if num, ok := issue["number"].(float64); ok && int(num) > info.LastIssue {
//				info.LastIssue = int(num)
//			}
//		}
//	}
//
//	// Последний pull request
//	var prs []map[string]interface{}
//	if err := GetJSON(fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls", owner, repo), &prs); err == nil {
//		for _, pr := range prs {
//			if num, ok := pr["number"].(float64); ok && int(num) > info.LastPR {
//				info.LastPR = int(num)
//			}
//		}
//	}
//
//	return info, nil
//}
