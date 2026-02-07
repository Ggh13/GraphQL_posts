package comment_repository

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

func (r Repository) Create(ctx context.Context, Comment *model.Comment) (*model.Comment, error) {
	Comment.ID = int32(len(r.localstorage.Comments)) + 1

	if int(Comment.User.ID) > len(r.localstorage.Users) { // Проверяем существование юзера
		errorW := fmt.Sprint("CommentRepository.Create: User author does not exist")
		logger.GetLoggerFromCtx(ctx).Info(ctx, errorW)
		return nil, fmt.Errorf(errorW)
	}

	if int(Comment.PostID) > len(r.localstorage.Posts) { // Проверяем существование поста
		errorW := fmt.Sprint("CommentRepository.Create: Post with the id does not exist")
		logger.GetLoggerFromCtx(ctx).Info(ctx, errorW)
		return nil, fmt.Errorf(errorW)
	}

	if Comment.ParentIDComment > 0 { // Если указан родитель комментария

		if int(Comment.ParentIDComment) <= len(r.localstorage.Comments) { // Проверяем его существование

			if Comment.PostID != r.localstorage.Comments[Comment.ParentIDComment-1].PostID { //А теперь проверяем что нет ошибки в указании поста куда пишется комментарий
				errorW := fmt.Sprint("CommentRepository.Create: PostId must be same like parent Comment PostID")
				logger.GetLoggerFromCtx(ctx).Info(ctx, errorW)
				return nil, fmt.Errorf(errorW)
			}

		} else {
			errorW := fmt.Sprint("CommentRepository.Create: Parent Id ( Parent comment ) must be exist")
			logger.GetLoggerFromCtx(ctx).Info(ctx, errorW)
			return nil, fmt.Errorf(errorW)
		}
	}

	Comment.User = &r.localstorage.Users[Comment.User.ID-1]

	r.localstorage.Comments = append(r.localstorage.Comments, *Comment)

	if Comment.ParentIDComment <= 0 { // Если родитель комментария не указан, то делаем его корневым для данного поста
		r.localstorage.Posts[Comment.PostID-1].Comments = append(r.localstorage.Posts[Comment.PostID-1].Comments, Comment)
		return Comment, nil
	}

	//Если родитель указан, то ищем родителя и добавляем в его дочерние комментарии
	r.localstorage.Comments[Comment.ParentIDComment-1].Comments = append(r.localstorage.Comments[Comment.ParentIDComment-1].Comments, Comment)

	return Comment, nil
}
func (r Repository) Update(ctx context.Context, Comment *model.Comment) (bool, error) {
	return false, nil
}
func (r Repository) Get(ctx context.Context, CommentId int) (*model.Comment, error) {
	if CommentId > len(r.localstorage.Comments) { // Проверяем существование комментария
		errorW := fmt.Sprint("CommentRepository.Get: User does not exist")
		logger.GetLoggerFromCtx(ctx).Info(ctx, errorW)
		return nil, fmt.Errorf(errorW)
	}
	return &r.localstorage.Comments[CommentId-1], nil
}
func (r Repository) Delete(ctx context.Context, Post int) (bool, error) {
	return false, nil
}
func (r Repository) GetAllCommentOfPost(ctx context.Context, PostId int) ([]*model.Comment, error) {
	if PostId > len(r.localstorage.Posts) || PostId < 0 {
		errorW := fmt.Sprint("CommentRepository.Get: User does not exist")
		logger.GetLoggerFromCtx(ctx).Info(ctx, errorW)
		return nil, fmt.Errorf(errorW)
	}
	var res []*model.Comment // Берем все комментариия соответсвующего поста
	for _, j := range r.localstorage.Comments {
		if j.PostID == int32(PostId) {
			res = append(res, &j)
		}
	}
	return res, nil
}

func (r Repository) GetPostId(ctx context.Context, CommentId int) (int, error) {

	var postId int
	if CommentId >= len(r.localstorage.Comments) { //Проверяем существование комментария
		postId = int(r.localstorage.Comments[CommentId-1].PostID)
	} else {
		errorW := fmt.Sprint("Failed to get post by ID comment %s %w", CommentId)
		logger.GetLoggerFromCtx(ctx).Info(ctx, errorW)
		return -1, fmt.Errorf(errorW)
	}

	return postId, nil
}
