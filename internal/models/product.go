package models

import "github.com/jackc/pgx/v5/pgtype"

// Product models
type Product struct {
	ID        int         `json:"id"`
	Name      string      `json:"name"`
	Price     int         `json:"price"`
	CreatedAt pgtype.Time `json:"created_at"`
}
