package sell

import (
	"context"

	"github.com/Lidne/praktika_MAI/internal/models"
)

// SellRepository Sell
type SellRepository interface {
	Create(ctx context.Context, product *models.Sell) error
	Update(ctx context.Context, product *models.Sell) error
	GetByID(ctx context.Context, id string) (*models.Sell, error)
	FindAll(ctx context.Context) ([]models.Sell, error)
	Delete(ctx context.Context, id int) error
}
