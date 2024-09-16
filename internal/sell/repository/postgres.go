package repository

import (
	"context"
	"fmt"
	"github.com/Lidne/praktika_MAI/internal/models"
	"github.com/Lidne/praktika_MAI/internal/sell"
	"github.com/Lidne/praktika_MAI/pkg/postgres"
)

// sellRepo
type sellRepo struct {
	client postgres.Client
}

// NewSellRepo sellRepo constructor
func NewSellRepo(client postgres.Client) sell.SellRepository {
	return &sellRepo{client: client}
}

func (r *sellRepo) Create(ctx context.Context, sell *models.Sell) error {
	q := `INSERT INTO bargains (user_id, product_id) VALUES ($1, $2) returning id`
	if err := r.client.QueryRow(ctx, q, sell.UserId, sell.ProductId).Scan(&sell.ID); err != nil {
		fmt.Errorf(fmt.Sprintf("SQL Error. Falied to create sell: %s", err.Error()))
		return err
	}

	return nil
}

func (r *sellRepo) Update(ctx context.Context, sell *models.Sell) error {
	q := `UPDATE bargains SET user_id=$1, product_id=$2 WHERE id=$3`
	_, err := r.client.Exec(ctx, q, &sell.UserId, &sell.ProductId, &sell.ID)
	if err != nil {
		fmt.Errorf("SQL Error. Failed to update sell: %s", err.Error())
		return err
	}
	return nil
}

func (r *sellRepo) GetByID(ctx context.Context, id string) (*models.Sell, error) {
	q := `SELECT id, user_id, product_id, updatedat FROM bargains WHERE id=$1`
	sell := &models.Sell{}
	if err := r.client.QueryRow(ctx, q, id).Scan(&sell.ID, &sell.UserId, &sell.ProductId, &sell.UpdatedAt); err != nil {
		fmt.Errorf("SQL Error. Failed to get sell by ID: %s", err.Error())
		return nil, err
	}
	return sell, nil
}

func (r *sellRepo) FindAll(ctx context.Context) ([]models.Sell, error) {
	q := `SELECT id, user_id, product_id, updatedat FROM bargains`
	rows, err := r.client.Query(ctx, q)
	if err != nil {
		fmt.Errorf("SQL Error. Failed to find all sells: %s", err.Error())
		return nil, err
	}
	defer rows.Close()

	sells := []models.Sell{}
	for rows.Next() {
		product := models.Sell{}
		if err := rows.Scan(&product.ID, &product.UserId, &product.ProductId, &product.UpdatedAt); err != nil {
			fmt.Errorf("SQL Error. Failed to scan sell: %s", err.Error())
			return nil, err
		}
		sells = append(sells, product)
	}

	return sells, nil
}

func (r *sellRepo) Delete(ctx context.Context, id int) error {
	q := `DELETE FROM bargains WHERE id=$1`
	_, err := r.client.Exec(ctx, q, id)
	if err != nil {
		fmt.Errorf("SQL Error. Failed to delete sell: %s", err.Error())
		return err
	}
	return nil
}

func (r *sellRepo) SelectByTime(ctx context.Context, time string) ([]models.Sell, error) {
	q := `SELECT * FROM bargains WHERE updatedat >= NOW() - INTERVAL $1`
	rows, err := r.client.Query(ctx, q, time)
	if err != nil {
		fmt.Errorf("SQL Error. Failed to find all sells: %s", err.Error())
		return nil, err
	}
	defer rows.Close()

	sells := []models.Sell{}
	for rows.Next() {
		sell := models.Sell{}
		if err := rows.Scan(&sell.ID, &sell.UserId, &sell.ProductId, &sell.UpdatedAt); err != nil {
			fmt.Errorf("SQL Error. Failed to scan sell: %s", err.Error())
			return nil, err
		}
		sells = append(sells, sell)
	}

	return sells, nil
}
