package comment_repository

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

func (r Repository) Create(ctx context.Context, Comment *model.Comment) (*model.Comment, error) {
	Comment.ID = int32(len(r.localstorage.Comments)) + 1

	if int(Comment.User.ID) > len(r.localstorage.Users) {
		return nil, fmt.Errorf("User author does not exist")
	}

	if int(Comment.PostID) > len(r.localstorage.Posts) {
		return nil, fmt.Errorf("Post with the id does not exist")
	}
	if Comment.ParentIDComment > 0 {

		if int(Comment.ParentIDComment) <= len(r.localstorage.Comments) {

			if Comment.PostID != r.localstorage.Comments[Comment.ParentIDComment-1].PostID {
				return nil, fmt.Errorf("PostId must be same like parent Comment PostID")
			}

		} else {
			return nil, fmt.Errorf("Parent Id ( Parent comment ) must be exist")
		}
	}

	Comment.User = &r.localstorage.Users[Comment.User.ID-1]

	r.localstorage.Comments = append(r.localstorage.Comments, *Comment)

	if Comment.ParentIDComment < 0 {
		r.localstorage.Posts[Comment.PostID-1].Comments = append(r.localstorage.Posts[Comment.PostID-1].Comments, Comment)
		return Comment, nil
	}

	r.localstorage.Comments[Comment.ParentIDComment-1].Comments = append(r.localstorage.Comments[Comment.ParentIDComment-1].Comments, Comment)

	return Comment, nil
}
func (r Repository) Update(ctx context.Context, Comment *model.Comment) (bool, error) {
	return false, nil
}
func (r Repository) Get(ctx context.Context, CommentId int) (*model.Comment, error) {
	if CommentId > len(r.localstorage.Comments) {
		return nil, fmt.Errorf("User does not exist")
	}
	return &r.localstorage.Comments[CommentId-1], nil
}
func (r Repository) Delete(ctx context.Context, Post int) (bool, error) {
	return false, nil
}
func (r Repository) GetAllPost(ctx context.Context, PostId int) ([]*model.Comment, error) {
	if PostId > len(r.localstorage.Posts) || PostId < 0 {
		return nil, fmt.Errorf("User does not exist")
	}
	return r.localstorage.Posts[PostId-1].Comments, nil
}

func (r Repository) GetPostId(ctx context.Context, CommentId int) (int, error) {

	var postId int
	if CommentId < len(r.localstorage.Comments) {
		postId = int(r.localstorage.Comments[CommentId-1].PostID)
	} else {
		return -1, fmt.Errorf("Failed to get post by ID comment %s %w", CommentId)
	}

	return postId, nil
}
