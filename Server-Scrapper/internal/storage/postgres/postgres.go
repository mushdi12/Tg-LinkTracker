package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	. "server-scrapper/internal/storage"
)

type StorageInterface interface {
	InsertUser(ctx context.Context, chatID int64, username string) bool
	InsertLink(ctx context.Context, chatID int64, linkUrl string, linkCategory string) bool
	DeleteLink(ctx context.Context, chatID int64, linkUrl string) (bool, string)
	UserExist(ctx context.Context, chatID int64) bool
	GetUserLinks(ctx context.Context, chatID int64) UserLinks
}

type Storage struct {
	POSTGRES *pgxpool.Pool
}

func New(postrgresData string) (*Storage, error) {
	const op = "storage.postgres.ConnectBD"
	pool, err := pgxpool.New(context.Background(), postrgresData)
	if err != nil {
		log.Fatal(op, err)
		//return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &Storage{POSTGRES: pool}, nil
}

func (storage *Storage) InsertUser(ctx context.Context, chatID int64, username string) bool {
	const op = "storage.postgres.Insert"
	_, err := storage.POSTGRES.Exec(ctx, "INSERT INTO USERS (chat_id,username) VALUES ($1,$2)",
		chatID, username)
	if err != nil {
		fmt.Errorf("%s: %w", op, err)
		return false
	}
	return true
}
func (storage *Storage) InsertLink(ctx context.Context, chatID int64, linkUrl string, linkCategory string) bool {
	const op = "storage.postgres.Insert"
	_, err := storage.POSTGRES.Exec(ctx, "INSERT INTO USERS_LINKS (chat_id,irl,category) VALUES ($1,$2,$3)",
		chatID, linkUrl, linkCategory)
	if err != nil {
		fmt.Errorf("%s: %w", op, err)
		return false
	}
	return true
}

func (storage *Storage) DeleteLink(ctx context.Context, chatID int64, linkUrl string) (bool, string) {
	const op = "storage.postgres.DeleteLink"
	answer, err := storage.POSTGRES.Exec(ctx, "DELETE FROM USERS_LINKS WHERE chat_id = $1 AND irl = $2",
		chatID, linkUrl)
	if err != nil {
		fmt.Errorf("%s: %w", op, err)
		return false, "" // произошла ошибка в бд
	}
	if answer.RowsAffected() == 0 {
		return true, "У пользователя нет ссылки"
	}

	return true, "Ссылка удалена"
}

func (storage *Storage) UserExist(ctx context.Context, chatID int64) bool {
	const op = "storage.postgres.DeleteLink"

	var count int
	err := storage.POSTGRES.QueryRow(ctx, "SELECT COUNT(*) FROM users WHERE chat_id=$1", chatID).Scan(&count)
	if err != nil {
		fmt.Errorf("%s: %w", op, err)
	}
	return count > 0
}

func (storage *Storage) GetUserLinks(ctx context.Context, chatID int64) UserLinks {
	const op = "storage.postgres.GetUserLinks"

	rows, err := storage.POSTGRES.Query(ctx, "SELECT irl, category FROM USERS_LINKS WHERE chat_id = $1", chatID)
	if err != nil {
		fmt.Errorf("%s: %w", op, err)
		return UserLinks{}
	}
	defer rows.Close()

	var links UserLinks

	for rows.Next() {
		var url, category string
		if err := rows.Scan(&url, &category); err != nil {
			fmt.Errorf("%s: %w", op, err)
			return UserLinks{}
		}
		link := UserLink{LinkUrl: url, LinkCategory: category}
		links = append(links, link)

	}

	if err := rows.Err(); err != nil {
		fmt.Errorf("%s: %w", op, err)
		return UserLinks{}
	}
	return links
}
