package comment_service

import (
	"context"
	"fmt"
	"qraphQL_posts/api/graph"
	"qraphQL_posts/api/graph/model"
	"qraphQL_posts/pkg/logger"
	"unicode/utf8"
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
	GetPostId(ctx context.Context, CommentId int) (int, error)
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

	idiPost := Comment.PostID

	if utf8.RuneCountInString(Comment.Content) > 2000 {
		return nil, fmt.Errorf("Your comment len %v. Maximum acepted len is 2000", utf8.RuneCountInString(Comment.Content))
	}

	PostToComment, err := s.postService.Get(ctx, int(idiPost))

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

func (s Service) Get(ctx context.Context, CommentID int) ([]*model.Comment, error) {
	postId, err := s.repo.GetPostId(ctx, CommentID)
	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}
	comments, err := s.repo.GetAllPost(ctx, postId)
	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}

	var commentbyId []*model.Comment
	for _, j := range comments {
		if j.ID == int32(CommentID) {

			commentbyId = append(commentbyId, j)
			var queue []*model.Comment
			queue = FindKids(int(j.ID), comments)

			for len(queue) != 0 {
				tar := queue[len(queue)-1]
				queue = queue[:len(queue)-1]
				commentbyId = append(commentbyId, tar)

				parentTar := FindKids(int(tar.ID), comments)

				queue = append(queue, parentTar...)
			}

		}

	}

	return commentbyId, nil
}

func (s Service) GetAllPost(ctx context.Context, PostId int) ([]*model.Comment, error) {
	comments, err := s.repo.GetAllPost(ctx, PostId)
	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}

	logger.GetLoggerFromCtx(ctx).Info(ctx, fmt.Sprintf("there len %v"))

	var preorityComments []*model.Comment
	for _, j := range comments {
		if !CheckSlice(j, preorityComments) && (j.ParentIDComment < 1) {

			preorityComments = append(preorityComments, j)
			var queue []*model.Comment
			queue = FindKids(int(j.ID), comments)
			logger.GetLoggerFromCtx(ctx).Info(ctx, fmt.Sprintf("there len Kids %v", len(queue)))

			for len(queue) != 0 {
				tar := queue[len(queue)-1]
				queue = queue[:len(queue)-1]
				preorityComments = append(preorityComments, tar)

				logger.GetLoggerFromCtx(ctx).Info(ctx, fmt.Sprintf("Now watch id %v", tar.ID))

				parentTar := FindKids(int(tar.ID), comments)

				queue = append(queue, parentTar...)
			}
			logger.GetLoggerFromCtx(ctx).Info(ctx, fmt.Sprintf("there len preority ID %v", len(preorityComments)))

		}

	}

	return preorityComments, nil
}

func CheckSlice(commentTarget *model.Comment, comments []*model.Comment) bool {
	for _, j := range comments {
		if j != nil && commentTarget == j {
			return true
		}
	}
	return false
}
func FindKids(IdParent int, comments []*model.Comment) []*model.Comment {
	var res []*model.Comment
	for _, j := range comments {
		if j != nil && (IdParent == int(j.ParentIDComment)) {
			res = append(res, j)
		}
	}
	return res
}

/*
	//Second stage. Start write all comments preority

	//Third stage find Target Comment
	fl := false
	var res []*model.Comment
	for _, j := range preorityComments {
		if fl {
			if res[len(res)-1].ID != j.ParentIDComment {
				return res, nil
			}
		}
	}
	return comments, nil
*/
