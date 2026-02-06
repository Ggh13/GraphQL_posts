package post_repository

import (
	"context"
	"fmt"
	"qraphQL_posts/api/graph/model"
	localstorage "qraphQL_posts/pkg/localStorage"
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
	Post.ID = int32(len(r.localstorage.Posts)) + 1
	if int(Post.User.ID) > len(r.localstorage.Users) {
		return nil, fmt.Errorf("User author does not exist")
	}
	Post.User = &r.localstorage.Users[Post.User.ID-1]
	r.localstorage.Posts = append(r.localstorage.Posts, *Post)

	return Post, nil
}

func (r Repository) Update(ctx context.Context, Post *model.Post) (bool, error) {
	idi := Post.ID
	if idi > int32(len(r.localstorage.Posts)) {
		return false, fmt.Errorf("The id post does not exist")
	}
	r.localstorage.Posts[idi-1].Commentable = *&Post.Commentable
	return true, nil
}
func (r Repository) Get(ctx context.Context, PostId int) (*model.Post, error) {
	if PostId > len(r.localstorage.Posts) {
		return nil, fmt.Errorf("User does not exist")
	}
	return &r.localstorage.Posts[PostId-1], nil
}
func (r Repository) Delete(ctx context.Context, Post int) (bool, error) {
	return false, nil
}

func (r Repository) GetAllPosts(ctx context.Context) ([]*model.Post, error) {
	posts := make([]*model.Post, len(r.localstorage.Posts))
	for i := range r.localstorage.Posts {
		posts[i] = &r.localstorage.Posts[i]
	}

	return posts, nil
}
