package user_psg_repository

import (
	"context"
	"fmt"
	"qraphQL_posts/api/graph/model"
	"qraphQL_posts/pkg/logger"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	queryCreateUser = `
        INSERT INTO users (name, surname) 
        VALUES ($1, $2)
        RETURNING id
    `

	queryGetUser = `
        SELECT id, name, surname
        FROM users 
        WHERE id = $1
    `
)

type Repository struct {
	pgDB *pgxpool.Pool
}

/*
	type Repository interface {
		Get(ctx context.Context) (bool, error)
		Create(ctx context.Context, User *model.User) (bool, error)
		Update(ctx context.Context, User *model.User) (bool, error)
		Delete(ctx context.Context, User *model.User) (bool, error)
	}
*/

func New(ctx context.Context, pgDB *pgxpool.Pool) *Repository {
	return &Repository{pgDB: pgDB}
}

func (r Repository) Create(ctx context.Context, User *model.User) (*model.User, error) {

	var id string
	err := r.pgDB.QueryRow(ctx, queryCreateUser,
		User.Name,
		User.Surname,
	).Scan(&id)

	if err != nil {
		logger.GetLoggerFromCtx(ctx).Info(ctx, fmt.Sprint("UserRepository.Create: failed to create user: %w", err))
		return nil, fmt.Errorf("UserRepository.Create: failed to create user: %w", err)
	}
	idInt, err := strconv.Atoi(id)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Info(ctx, fmt.Sprint("UserRepository.Create: failed to create user: %w", err))
		return nil, fmt.Errorf("UserRepository.Create: failed to create user: %w", err)
	}
	User.ID = int32(idInt)
	return User, nil
}
func (r Repository) Update(ctx context.Context, User *model.User) (bool, error) {
	return false, nil
}
func (r Repository) Get(ctx context.Context, userID int) (*model.User, error) {

	var user model.User
	err := r.pgDB.QueryRow(ctx, queryGetUser, userID).Scan(
		&user.ID,
		&user.Name,
		&user.Surname,
	)

	if err != nil {
		logger.GetLoggerFromCtx(ctx).Info(ctx, fmt.Sprint("UserRepository.Create: Failed to get user by ID %s %w", userID, err))
		return nil, fmt.Errorf("UserRepository.Create: Failed to get user by ID %s %w", userID, err)
	}

	return &user, nil
}
func (r Repository) Delete(ctx context.Context, UserId int) (bool, error) {
	return false, nil
}
