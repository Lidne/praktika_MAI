package repository

import (
	"context"
	"fmt"
	"github.com/Lidne/praktika_MAI/internal/models"
	"github.com/Lidne/praktika_MAI/internal/product"

	"github.com/Lidne/praktika_MAI/pkg/postgres"
)

// userRepo
type userRepo struct {
	client postgres.Client
}

// NewUserRepo userRepo constructor
func NewUserRepo(client postgres.Client) product.ProductRepository {
	return &userRepo{client: client}
}

func (r *userRepo) Create(ctx context.Context, product *models.Product) error {
	q := `INSERT INTO users (name, password, ) VALUES ($1, $2) returning id`
	if err := r.client.QueryRow(ctx, q, product.Name, product.Price).Scan(&product.ID); err != nil {
		fmt.Errorf(fmt.Sprintf("SQL Error. Falied to create product: %s", err.Error()))
		return err
	}

	return nil
}

func (r *userRepo) Update(ctx context.Context, product *models.Product) error {
	q := `UPDATE users SET name=$1, price=$2 WHERE id=$3`
	_, err := r.client.Exec(ctx, q, &product.Name, &product.Price, &product.ID)
	if err != nil {
		fmt.Errorf("SQL Error. Failed to update product: %s", err.Error())
		return err
	}
	return nil
}

func (r *userRepo) GetByID(ctx context.Context, id int) (*models.Product, error) {
	q := `SELECT id, name, price, updatedat FROM products WHERE id=$1`
	product := &models.Product{}
	if err := r.client.QueryRow(ctx, q, id).Scan(&product.ID, &product.Name, &product.Price, &product.CreatedAt); err != nil {
		fmt.Errorf("SQL Error. Failed to get product by ID: %s", err.Error())
		return nil, err
	}
	return product, nil
}

func (r *userRepo) FindAll(ctx context.Context) ([]models.Product, error) {
	q := `SELECT id, name, price, updatedat FROM products`
	rows, err := r.client.Query(ctx, q)
	if err != nil {
		fmt.Errorf("SQL Error. Failed to find all products: %s", err.Error())
		return nil, err
	}
	defer rows.Close()

	products := []models.Product{}
	for rows.Next() {
		product := models.Product{}
		if err := rows.Scan(&product.ID, &product.Name, &product.Price, &product.CreatedAt); err != nil {
			fmt.Errorf("SQL Error. Failed to scan product: %s", err.Error())
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}

func (r *userRepo) Delete(ctx context.Context, id int) error {
	q := `DELETE FROM products WHERE id=$1`
	_, err := r.client.Exec(ctx, q, id)
	if err != nil {
		fmt.Errorf("SQL Error. Failed to delete product: %s", err.Error())
		return err
	}
	return nil
}
