package user_service

import (
	"context"
	"fmt"
	"qraphQL_posts/api/graph"
	"qraphQL_posts/api/graph/model"
	"strconv"
)

/*
	Create(ctx context.Context, Comment *model.Comment) (*model.Comment, error)
	Get(ctx context.Context, CommentId int) (*model.Comment, error)
	GetAllPost(ctx context.Context, PostId int) ([]*model.Comment, error)
*/

type Repository interface {
	Get(ctx context.Context, CommentId int) (*model.Comment, error)
	Create(ctx context.Context, Comment *model.Comment) (*model.Comment, error)
	Update(ctx context.Context, Comment *model.Comment) (bool, error)
	Delete(ctx context.Context, CommentId int) (bool, error)
	GetAllPost(ctx context.Context, PostId int) ([]*model.Comment, error)
}

type Service struct {
	repo        Repository
	postService graph.PostService
}

func New(ctx context.Context, r Repository, pS graph.PostService) Service {
	return Service{repo: r,
		postService: pS}
}

func (s Service) Create(ctx context.Context, Comment *model.Comment) (*model.Comment, error) {

	idiPost, err := strconv.Atoi(Comment.PostID)
	if err != nil {
		return nil, fmt.Errorf("Error in id Post")
	}

	PostToComment, err := s.postService.Get(ctx, idiPost)

	if err != nil {
		return nil, fmt.Errorf("There are not post with id %v", idiPost)
	}

	if !PostToComment.Commentable {
		return nil, fmt.Errorf("This post does not accept comments")
	}

	fl, err := s.repo.Create(ctx, Comment)
	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}

	return fl, nil
}

func (s Service) Get(ctx context.Context, CommentID int) (*model.Comment, error) {
	res, err := s.repo.Get(ctx, CommentID)
	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}
	return res, nil
}

func (s Service) GetAllPost(ctx context.Context, PostId int) ([]*model.Comment, error) {
	res, err := s.repo.GetAllPost(ctx, PostId)
	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}
	return res, nil
}
