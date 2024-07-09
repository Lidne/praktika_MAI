package repository

import (
	"context"
	"fmt"
	"github.com/Lidne/praktika_MAI/internal/models"
	"github.com/Lidne/praktika_MAI/internal/user"

	"github.com/Lidne/praktika_MAI/pkg/postgres"
)

// userRepo
type userRepo struct {
	client postgres.Client
}

// NewUserRepo userRepo constructor
func NewUserRepo(client postgres.Client) user.UserRepository {
	return &userRepo{client: client}
}

func (r *userRepo) Create(ctx context.Context, user *models.User) error {
	q := `INSERT INTO users (name, login, password) VALUES ($1, $2, $3) returning id`
	if err := r.client.QueryRow(ctx, q, user.Name, user.Login, user.Password).Scan(&user.ID); err != nil {
		fmt.Errorf(fmt.Sprintf("SQL Error. Falied to create user: %s", err.Error()))
		return err
	}

	return nil
}

func (r *userRepo) Update(ctx context.Context, user *models.User) error {
	q := `UPDATE users SET name=$1, login=$2, password=$3 WHERE id=$4`
	_, err := r.client.Exec(ctx, q, &user.Name, &user.Login, &user.Password, &user.ID)
	if err != nil {
		fmt.Errorf("SQL Error. Failed to update user: %s", err.Error())
		return err
	}
	return nil
}

func (r *userRepo) GetByID(ctx context.Context, id string) (*models.User, error) {
	q := `SELECT id, name, updatedat, login, password, isadmin FROM users WHERE id=$1`
	user := &models.User{}
	if err := r.client.QueryRow(ctx, q, id).Scan(&user.ID, &user.Name, &user.UpdatedAt, &user.Login, &user.Password, &user.IsAdmin); err != nil {
		fmt.Errorf("SQL Error. Failed to get user by ID: %s", err.Error())
		return nil, err
	}
	return user, nil
}

func (r *userRepo) FindAll(ctx context.Context) ([]models.User, error) {
	q := `SELECT id, name, updatedat, login, password, isadmin FROM users`
	rows, err := r.client.Query(ctx, q)
	if err != nil {
		fmt.Errorf("SQL Error. Failed to find all users: %s", err.Error())
		return nil, err
	}
	defer rows.Close()

	users := []models.User{}
	for rows.Next() {
		user := models.User{}
		if err := rows.Scan(&user.ID, &user.Name, &user.UpdatedAt, &user.Login, &user.Password, &user.IsAdmin); err != nil {
			fmt.Errorf("SQL Error. Failed to scan user: %s", err.Error())
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *userRepo) Delete(ctx context.Context, id int) error {
	q := `DELETE FROM users WHERE id=$1`
	_, err := r.client.Exec(ctx, q, id)
	if err != nil {
		fmt.Errorf("SQL Error. Failed to delete user: %s", err.Error())
		return err
	}
	return nil
}
