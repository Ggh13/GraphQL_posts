package user_repository

import (
	"context"
	"fmt"
	"qraphQL_posts/api/graph/model"
	localstorage "qraphQL_posts/pkg/localStorage"
	"strconv"
)

type Repository struct {
	localstorage *localstorage.Storage
}

/*
	type Repository interface {
		Get(ctx context.Context) (bool, error)
		Create(ctx context.Context, User *model.User) (bool, error)
		Update(ctx context.Context, User *model.User) (bool, error)
		Delete(ctx context.Context, User *model.User) (bool, error)
	}
*/
func New(localstorageR *localstorage.Storage) Repository {
	return Repository{
		localstorage: localstorageR,
	}
}

func (r Repository) Create(ctx context.Context, Post *model.Post) (*model.Post, error) {
	Post.ID = fmt.Sprint(len(r.localstorage.Posts))
	r.localstorage.Posts = append(r.localstorage.Posts, *Post)
	return Post, nil
}

func (r Repository) Update(ctx context.Context, Post *model.Post) (bool, error) {
	idi, err := strconv.Atoi(Post.ID)
	if err != nil {
		return false, fmt.Errorf("Error of ID Post")
	}
	r.localstorage.Posts[idi].Commentable = *&Post.Commentable
	return true, nil
}
func (r Repository) Get(ctx context.Context, PostId int) (*model.Post, error) {
	if PostId >= len(r.localstorage.Posts) {
		return nil, fmt.Errorf("User does not exist")
	}
	return &r.localstorage.Posts[PostId], nil
}
func (r Repository) Delete(ctx context.Context, Post int) (bool, error) {
	return false, nil
}
