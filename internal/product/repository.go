package product

import (
	"context"

	"github.com/Lidne/praktika_MAI/internal/models"
)

// ProductRepository Sell
type ProductRepository interface {
	Create(ctx context.Context, product *models.Product) error
	Update(ctx context.Context, product *models.Product) error
	GetByID(ctx context.Context, id string) (*models.Product, error)
	FindAll(ctx context.Context) ([]models.Product, error)
	Delete(ctx context.Context, id int) error
}
