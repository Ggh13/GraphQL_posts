package user_repository

import (
	"context"
	"fmt"
	"qraphQL_posts/api/graph/model"
	localstorage "qraphQL_posts/pkg/localStorage"
	"qraphQL_posts/pkg/logger"
)

type Repository struct {
	localstorage *localstorage.Storage
}

func New(localstorageR *localstorage.Storage) Repository {
	return Repository{
		localstorage: localstorageR,
	}
}

func (r Repository) Create(ctx context.Context, User *model.User) (*model.User, error) {
	User.ID = int32(len(r.localstorage.Users)) + 1
	r.localstorage.Users = append(r.localstorage.Users, *User)
	return User, nil
}
func (r Repository) Update(ctx context.Context, User *model.User) (bool, error) {
	return false, nil
}
func (r Repository) Get(ctx context.Context, userID int) (*model.User, error) {
	if userID > len(r.localstorage.Users) {
		logger.GetLoggerFromCtx(ctx).Info(ctx, fmt.Sprint("UserRepository.Create: User does not exist"))
		return nil, fmt.Errorf("UserRepository.Create: User does not exist")
	}
	return &r.localstorage.Users[userID-1], nil
}
func (r Repository) Delete(ctx context.Context, UserId int) (bool, error) {
	return false, nil
}
