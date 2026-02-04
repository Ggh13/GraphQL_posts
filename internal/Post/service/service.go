package user_service

import (
	"context"
	"fmt"
	"qraphQL_posts/api/graph/model"
)

/*
	Create(ctx context.Context, Post *model.Post) (bool, error)
	Update(ctx context.Context, Post *model.Post) (bool, error)
	Delete(ctx context.Context, postID int) (bool, error)
	Get(ctx context.Context, postID int) (*model.Post, error)
*/

type Repository interface {
	Get(ctx context.Context, PostId int) (*model.Post, error)
	Create(ctx context.Context, Post *model.Post) (*model.Post, error)
	Update(ctx context.Context, Post *model.Post) (bool, error)
	Delete(ctx context.Context, PostId int) (bool, error)
}

type Service struct {
	repo Repository
}

func New(ctx context.Context, r Repository) Service {
	return Service{repo: r}
}

func (s Service) Create(ctx context.Context, Post *model.Post) (*model.Post, error) {
	fl, err := s.repo.Create(ctx, Post)
	if err != nil {
		return fl, fmt.Errorf("%s", err)
	}
	return fl, nil
}

func (s Service) Update(ctx context.Context, Post *model.Post) (bool, error) {
	fl, err := s.repo.Update(ctx, Post)
	if err != nil {
		return fl, fmt.Errorf("%s", err)
	}
	return true, nil
}

func (s Service) Delete(ctx context.Context, postID int) (bool, error) {

	return true, nil
}
func (s Service) Get(ctx context.Context, postID int) (*model.Post, error) {
	res, err := s.repo.Get(ctx, postID)
	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}
	return res, nil
}
