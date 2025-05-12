package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"server-scrapper/internal/lib/repositories"
	"server-scrapper/internal/storage"
)

type StorageInterface interface {
	InsertUser(ctx context.Context, chatID int64, username string) error
	InsertLink(ctx context.Context, chatID int64, linkUrl string, linkCategory string) error
	DeleteLink(ctx context.Context, chatID int64, linkUrl string) error
	InsertRepo(ctx context.Context, chatID int64, repoOwner, repoName, lastSHA string, lastIssue, lastPR int) error
	UserExist(ctx context.Context, chatID int64) (bool, error)
	GetUserLinks(ctx context.Context, chatID int64) (map[string]string, error)
	GetOldestRepos(ctx context.Context, count int) ([]repositories.WatchedRepo, error)
	UpdateRepoCheckTime(ctx context.Context, repoId int64) error
	UpdateRepo(ctx context.Context, repo repositories.WatchedRepo) error
}

type Storage struct {
	POSTGRES *pgxpool.Pool
}

func New(configData string) *Storage {
	const op = "storage.postgres.New"
	pool, err := pgxpool.New(context.Background(), configData)
	if err != nil {
		log.Fatalf("%s: %s", op, err)
	}
	return &Storage{POSTGRES: pool}
}

// User's methods
func (st *Storage) InsertUser(ctx context.Context, chatID int64, username string) error {
	const op = "storage.postgres.InsertUser"

	_, err := st.POSTGRES.Exec(ctx, "INSERT INTO USERS (chat_id,username) VALUES ($1,$2)", chatID, username)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == storage.UniqueViolationCode {
			return storage.ErrUserExists
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (st *Storage) UserExist(ctx context.Context, chatID int64) (bool, error) {
	const op = "storage.postgres.UserExist"

	var count int
	err := st.POSTGRES.QueryRow(ctx, "SELECT COUNT(*) FROM users WHERE chat_id=$1", chatID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	if count == 0 {
		return false, storage.ErrUserNotFound
	}
	return true, nil
}

// Link's methods
type UserLinks map[string]string

func (st *Storage) GetUserLinks(ctx context.Context, chatID int64) (map[string]string, error) {
	const op = "storage.postgres.GetUserLinks"

	rows, err := st.POSTGRES.Query(ctx, "SELECT url, category FROM USERS_LINKS WHERE chat_id = $1", chatID)
	defer rows.Close()
	if err != nil {
		_ = fmt.Errorf("%s: %w", op, err)
		return UserLinks{}, fmt.Errorf("%s: %w", op, err)
	}

	links := make(UserLinks)

	for rows.Next() {
		var url, category string
		if err := rows.Scan(&url, &category); err != nil {
			_ = fmt.Errorf("%s: %w", op, err)
			return UserLinks{}, fmt.Errorf("%s: %w", op, err)
		}
		links[url] = category

	}

	if err := rows.Err(); err != nil {
		_ = fmt.Errorf("%s: %w", op, err)
		return UserLinks{}, fmt.Errorf("%s: %w", op, err)
	}

	if len(links) == 0 {
		return nil, storage.ErrNoLinksFound
	}

	return links, nil
}

func (st *Storage) DeleteLink(ctx context.Context, chatID int64, linkUrl string) error {
	const op = "storage.postgres.DeleteLink"

	res, err := st.POSTGRES.Exec(ctx, "DELETE FROM USERS_LINKS WHERE chat_id = $1 AND url = $2", chatID, linkUrl)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if res.RowsAffected() == 0 {
		return storage.ErrLinkNotFound
	}

	return nil
}

func (st *Storage) InsertLink(ctx context.Context, chatID int64, linkUrl string, linkCategory string) error {
	const op = "storage.postgres.InsertLink"

	_, err := st.POSTGRES.Exec(ctx,
		`INSERT INTO USERS_LINKS (chat_id,url,category) 
	VALUES ($1,$2,$3)`, chatID, linkUrl, linkCategory)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == storage.UniqueViolationCode {
			return fmt.Errorf("%s: %w", op, storage.ErrLinkExists)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// Repo's methods
func (st *Storage) InsertRepo(ctx context.Context, chatID int64, repoOwner, repoName, lastSHA string, lastIssue, lastPR int) error {
	const op = "storage.postgres.InsertRepo"

	_, err := st.POSTGRES.Exec(ctx,
		`INSERT INTO WATCHED_REPOS (chat_id, repo_owner, repo_name, last_sha, last_issue, last_pr)
	VALUES ($1, $2, $3, $4, $5, $6)`, chatID, repoOwner, repoName, lastSHA, lastIssue, lastPR)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == storage.UniqueViolationCode {
			return fmt.Errorf("%s: %w", op, storage.ErrLinkExists)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (st *Storage) GetOldestRepos(ctx context.Context, count int) ([]repositories.WatchedRepo, error) {
	const op = "storage.postgres.GetOldestRepos"
	var repos []repositories.WatchedRepo

	query := `
		SELECT id, chat_id, repo_owner, repo_name, last_sha, last_issue, last_pr, updated_at
		FROM watched_repos
		ORDER BY updated_at DESC
		LIMIT $1;
	`

	rows, err := st.POSTGRES.Query(ctx, query, count)
	defer rows.Close()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	for rows.Next() {
		var repo repositories.WatchedRepo
		err := rows.Scan(&repo.ID, &repo.ChatID, &repo.RepoOwner, &repo.RepoName,
			&repo.LastSHA, &repo.LastIssue, &repo.LastPR, &repo.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("%s: ошибка при сканировании строки: %w", op, err)
		}

		repos = append(repos, repo)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: ошибка при обработке строк: %w", op, err)
	}

	return repos, nil
}

func (st *Storage) UpdateRepoCheckTime(ctx context.Context, repoId int64) error {
	const op = "storage.postgres.UpdateRepoCheckTime"
	_, err := st.POSTGRES.Exec(ctx, `UPDATE watched_repos SET updated_at = NOW() WHERE id = $1`, repoId)
	if err != nil {
		fmt.Printf("%s: %w", op, err)
	}
	return err
}

func (st *Storage) UpdateRepo(ctx context.Context, repo repositories.WatchedRepo) error {
	const op = "storage.postgres.UpdateRepo"

	_, err := st.POSTGRES.Exec(ctx, `UPDATE watched_repos 
		SET last_sha = $1, last_issue = $2, last_pr = $3, updated_at = NOW() WHERE id = $4
	`, repo.LastSHA, repo.LastIssue, repo.LastPR, repo.ID)

	if err != nil {
		fmt.Printf("%s: %w", op, err)
	}
	return err
}
