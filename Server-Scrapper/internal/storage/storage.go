package storage

import (
	"errors"
)

type Context struct {
	ChatID   int64
	UserName string
}

const UniqueViolationCode = "23505"

var (
	ErrLinkExists   = errors.New("link already exists")
	ErrUserExists   = errors.New("user already exists")
	ErrLinkNotFound = errors.New("link not found")
	ErrNoLinksFound = errors.New("link not found for this user")
	ErrUserNotFound = errors.New("user not found")
)
