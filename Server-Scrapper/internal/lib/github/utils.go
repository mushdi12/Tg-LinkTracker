package github

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

func parseRepoURL(rawURL string) (owner, repo string, err error) {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return "", "", err
	}
	parts := strings.Split(strings.Trim(parsed.Path, "/"), "/")
	if len(parts) < 2 {
		return "", "", fmt.Errorf("невалидная ссылка на репозиторий")
	}
	return parts[0], parts[1], nil
}

func parseCommit(data []byte) (string, string) {
	var commits []map[string]interface{}
	if err := json.Unmarshal(data, &commits); err != nil || len(commits) == 0 {
		return "", ""
	}

	sha, ok := commits[0]["sha"].(string)
	if !ok {
		return "", ""
	}

	commitInfo, ok := commits[0]["commit"].(map[string]interface{})
	if !ok {
		return "", ""
	}

	message, ok := commitInfo["message"].(string)
	if !ok {
		return "", ""
	}
	return sha, message
}

func parseIssue(data []byte) []map[string]interface{} {
	var issues []map[string]interface{}
	if err := json.Unmarshal(data, &issues); err != nil {
		return nil
	}
	return issues
}

func parsePR(data []byte) []map[string]interface{} {
	var prs []map[string]interface{}
	if err := json.Unmarshal(data, &prs); err != nil {
		return nil
	}
	return prs
}
