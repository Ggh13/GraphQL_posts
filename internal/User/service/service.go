package user_service

import (
	"context"
	"fmt"
	"qraphQL_posts/api/graph/model"
	"qraphQL_posts/pkg/logger"
)

type Repository interface {
	Get(ctx context.Context, UserId int) (*model.User, error)
	Create(ctx context.Context, User *model.User) (*model.User, error)
	Update(ctx context.Context, User *model.User) (bool, error)
	Delete(ctx context.Context, UserId int) (bool, error)
}

type Service struct {
	repo Repository
}

func New(ctx context.Context, r Repository) Service {
	return Service{repo: r}
}

func (s Service) Create(ctx context.Context, User *model.User) (*model.User, error) {
	fl, err := s.repo.Create(ctx, User)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Info(ctx, fmt.Sprintf("UserService.Create: %w", err))
		return fl, fmt.Errorf("UserService.Create: %w", err)
	}
	return fl, nil
}

func (s Service) Update(ctx context.Context, User *model.User) (bool, error) {

	return true, nil
}

func (s Service) Delete(ctx context.Context, userID int) (bool, error) {

	return true, nil
}
func (s Service) Get(ctx context.Context, userID int) (*model.User, error) {
	res, err := s.repo.Get(ctx, userID)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Info(ctx, fmt.Sprintf("UserService.Get: %w", err))
		return nil, fmt.Errorf("UserService.Get: %w", err)
	}
	return res, nil
}
