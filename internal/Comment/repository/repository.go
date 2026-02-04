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

func (r Repository) Create(ctx context.Context, Comment *model.Comment) (*model.Comment, error) {
	Comment.ID = fmt.Sprint(len(r.localstorage.Comments))
	r.localstorage.Comments = append(r.localstorage.Comments, *Comment)

	idiPost, err := strconv.Atoi(Comment.PostID)
	if err != nil {
		return nil, fmt.Errorf("Error in id Post")
	}

	idParentComment, err := strconv.Atoi(Comment.ParentIDComment)
	if err != nil {
		return nil, fmt.Errorf("Error in id Post")
	}

	if idParentComment < 0 {
		r.localstorage.Posts[idiPost].Comments = append(r.localstorage.Posts[idiPost].Comments, Comment)
		return Comment, nil
	}

	Comment_tree := r.localstorage.Posts[idiPost].Comments
	queue := []*model.Comment{}
	for _, i := range Comment_tree {
		queue = append(queue, i)
	}

	for len(queue) >= 1 {
		tar := queue[len(queue)-1]
		queue = queue[:len(queue)-1]
		if tar.ID == Comment.ParentIDComment {
			tar.Comments = append(tar.Comments, Comment)
			return Comment, nil
		}
		for _, i := range tar.Comments {
			queue = append(queue, i)
		}
	}

	return nil, fmt.Errorf("There error with ParentId comment")
}
func (r Repository) Update(ctx context.Context, Comment *model.Comment) (bool, error) {
	return false, nil
}
func (r Repository) Get(ctx context.Context, CommentId int) (*model.Comment, error) {
	if CommentId >= len(r.localstorage.Comments) {
		return nil, fmt.Errorf("User does not exist")
	}
	return &r.localstorage.Comments[CommentId], nil
}
func (r Repository) Delete(ctx context.Context, Post int) (bool, error) {
	return false, nil
}
func (r Repository) GetAllPost(ctx context.Context, PostId int) ([]*model.Comment, error) {
	if PostId >= len(r.localstorage.Posts) || PostId < 0 {
		return nil, fmt.Errorf("User does not exist")
	}
	return r.localstorage.Posts[PostId].Comments, nil
}
