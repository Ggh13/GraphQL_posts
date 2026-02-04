package graph

import (
	"context"
	"qraphQL_posts/api/graph/model"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require
// here.

type UserService interface {
	Create(ctx context.Context, User *model.User) (*model.User, error)
	Update(ctx context.Context, User *model.User) (bool, error)
	Delete(ctx context.Context, userID int) (bool, error)
	Get(ctx context.Context, userID int) (*model.User, error)
}

type CommentService interface {
	Create(ctx context.Context, Comment *model.Comment) (*model.Comment, error)
	Get(ctx context.Context, CommentId int) (*model.Comment, error)
	GetAllPost(ctx context.Context, PostId int) ([]*model.Comment, error)
	//Delete(ctx context.Context, userID int) (bool, error)
}

type PostService interface {
	Create(ctx context.Context, Post *model.Post) (*model.Post, error)
	Update(ctx context.Context, Post *model.Post) (bool, error)
	Delete(ctx context.Context, postID int) (bool, error)
	Get(ctx context.Context, postID int) (*model.Post, error)
}

type Resolver struct {
	ctx context.Context

	userService    UserService
	postService    PostService
	commentService CommentService
}

func (*Resolver) NewResolver(ct context.Context, us UserService, ps PostService, cs CommentService) *Resolver {
	return &Resolver{ctx: ct,
		userService:    us,
		postService:    ps,
		commentService: cs,
	}
}
