package models

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type User struct {
	ID        int
	Name      string
	UpdatedAt pgtype.Timestamp
	Login     string
	Password  string
	IsAdmin   bool
}
