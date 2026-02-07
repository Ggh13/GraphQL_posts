package post_service

import (
	"context"
	"fmt"
	"qraphQL_posts/api/graph"
	"qraphQL_posts/api/graph/model"
	"qraphQL_posts/pkg/logger"
)

type Repository interface {
	Get(ctx context.Context, PostId int) (*model.Post, error)
	GetAllPosts(ctx context.Context) ([]*model.Post, error)
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
		logger.GetLoggerFromCtx(ctx).Info(ctx, fmt.Sprint("PostService.Create: %s", err))
		return fl, fmt.Errorf("PostService.Create: %s", err)
	}
	return fl, nil
}

func (s Service) Update(ctx context.Context, Post *model.Post) (bool, error) {
	//В функции реализован Update ТОЛЬКО для изменения прав на комментирование, для безопасности
	fl, err := s.repo.Update(ctx, Post)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Info(ctx, fmt.Sprint("PostService.Update: %s", err))
		return fl, fmt.Errorf("PostService.Update: %s", err)
	}
	return true, nil
}

func (s Service) Delete(ctx context.Context, postID int) (bool, error) {

	return true, nil
}
func (s Service) Get(ctx context.Context, postID int) (*model.Post, error) {
	res, err := s.repo.Get(ctx, postID)
	if err != nil {
		return nil, fmt.Errorf("PostService.Get: %w", err)
	}

	return res, nil
}

func (s Service) GetAllPost(ctx context.Context, limit int, offset int) ([]*model.Post, error) {
	posts, err := s.repo.GetAllPosts(ctx)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Info(ctx, fmt.Sprint("PostService.GetAllPost: Error. Cant get all post %2", err))
		return nil, fmt.Errorf("PostService.GetAllPost: Error. Cant get all post %2", err)
	}
	err = graph.CheckLimit(posts, &limit, &offset)
	if err != nil {
		return nil, fmt.Errorf("PostService.GetAllPost: %w", err)
	}
	return posts[offset : offset+limit], nil

}
