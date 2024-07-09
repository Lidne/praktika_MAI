package models

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type User struct {
	id        int
	name      string
	updatedAt pgtype.Timestamp
	login     string
	password  string
	isAdmin   bool
}
