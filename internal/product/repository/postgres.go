package repository

import (
	"context"
	"fmt"
	"github.com/Lidne/praktika_MAI/internal/models"
	"github.com/Lidne/praktika_MAI/internal/product"

	"github.com/Lidne/praktika_MAI/pkg/postgres"
)

// productRepo
type productRepo struct {
	client postgres.Client
}

// NewProductRepo productRepo constructor
func NewProductRepo(client postgres.Client) product.ProductRepository {
	return &productRepo{client: client}
}

func (r *productRepo) Create(ctx context.Context, product *models.Product) error {
	q := `INSERT INTO products (name, price) VALUES ($1, $2) returning id`
	if err := r.client.QueryRow(ctx, q, product.Name, product.Price).Scan(&product.ID); err != nil {
		fmt.Errorf(fmt.Sprintf("SQL Error. Falied to create product: %s", err.Error()))
		return err
	}

	return nil
}

func (r *productRepo) Update(ctx context.Context, product *models.Product) error {
	q := `UPDATE products SET name=$1, price=$2 WHERE id=$3`
	_, err := r.client.Exec(ctx, q, &product.Name, &product.Price, &product.ID)
	if err != nil {
		fmt.Errorf("SQL Error. Failed to update product: %s", err.Error())
		return err
	}
	return nil
}

func (r *productRepo) GetByID(ctx context.Context, id string) (*models.Product, error) {
	q := `SELECT id, name, price, updatedat FROM products WHERE id=$1`
	product := &models.Product{}
	if err := r.client.QueryRow(ctx, q, id).Scan(&product.ID, &product.Name, &product.Price, &product.CreatedAt); err != nil {
		fmt.Errorf("SQL Error. Failed to get product by ID: %s", err.Error())
		return nil, err
	}
	return product, nil
}

func (r *productRepo) FindAll(ctx context.Context) ([]models.Product, error) {
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

func (r *productRepo) Delete(ctx context.Context, id int) error {
	q := `DELETE FROM products WHERE id=$1`
	_, err := r.client.Exec(ctx, q, id)
	if err != nil {
		fmt.Errorf("SQL Error. Failed to delete product: %s", err.Error())
		return err
	}
	return nil
}
