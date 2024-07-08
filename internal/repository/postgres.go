package repository

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	productsDB         = "products"
	productsCollection = "products"
)

// postgresRepo
type postgresRepo struct {
	dbpool *pgxpool.Pool
}

// NewPostgresRepo postgresRepo constructor
func NewPostgresRepo(dbpool *pgxpool.Pool) *postgresRepo {
	return &postgresRepo{dbpool: dbpool}
}
