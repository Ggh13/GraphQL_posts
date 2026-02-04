package user_repository

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

func (r Repository) Create(ctx context.Context, User *model.User) (*model.User, error) {
	User.ID = fmt.Sprint(len(r.localstorage.Users))
	r.localstorage.Users = append(r.localstorage.Users, *User)
	return User, nil
}
func (r Repository) Update(ctx context.Context, User *model.User) (bool, error) {
	return false, nil
}
func (r Repository) Get(ctx context.Context, userID int) (*model.User, error) {
	if userID >= len(r.localstorage.Users) {
		return nil, fmt.Errorf("User does not exist")
	}
	return &r.localstorage.Users[userID], nil
}
func (r Repository) Delete(ctx context.Context, UserId int) (bool, error) {
	return false, nil
}
