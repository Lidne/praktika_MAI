package models

import "github.com/jackc/pgx/v5/pgtype"

type Sell struct {
	ID        int
	UserId    int
	ProductId int
	UpdatedAt pgtype.Timestamp
}
