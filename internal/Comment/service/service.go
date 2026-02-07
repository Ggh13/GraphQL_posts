package comment_service

import (
	"context"
	"fmt"
	"qraphQL_posts/api/graph"
	"qraphQL_posts/api/graph/model"
	"qraphQL_posts/pkg/logger"
	"unicode/utf8"
)

type Repository interface {
	Get(ctx context.Context, CommentId int) (*model.Comment, error)
	Create(ctx context.Context, Comment *model.Comment) (*model.Comment, error)
	Update(ctx context.Context, Comment *model.Comment) (bool, error)
	Delete(ctx context.Context, CommentId int) (bool, error)
	GetAllCommentOfPost(ctx context.Context, PostId int) ([]*model.Comment, error)
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
		errorW := fmt.Sprint("CommentService.Create: Your comment len %v. Maximum acepted len is 2000", utf8.RuneCountInString(Comment.Content))
		logger.GetLoggerFromCtx(ctx).Info(ctx, errorW)
		return nil, fmt.Errorf(errorW)
	}

	PostToComment, err := s.postService.Get(ctx, int(idiPost)) // Проверка существования поста с данным ID
	if err != nil {
		errorW := fmt.Sprint("CommentService.Create: There are not post with id %v", idiPost)
		logger.GetLoggerFromCtx(ctx).Info(ctx, errorW)
		return nil, fmt.Errorf(errorW)

	}

	if !PostToComment.Commentable { // Проверка что данный пост можно комментировать
		errorW := fmt.Sprint("CommentService.Create: This post does not accept comments")
		logger.GetLoggerFromCtx(ctx).Info(ctx, errorW)
		return nil, fmt.Errorf(errorW)
	}
	fl, err := s.repo.Create(ctx, Comment)
	if err != nil {
		errorW := fmt.Sprint("CommentService.Create: %s", err)
		logger.GetLoggerFromCtx(ctx).Info(ctx, errorW)
		return nil, fmt.Errorf(errorW)
	}

	return fl, nil
}

func (s Service) Get(ctx context.Context, CommentID int, limit int, offset int) ([]*model.Comment, error) {
	//Функция поиска конкретного коментария и всех его дочерних
	postId, err := s.repo.GetPostId(ctx, CommentID) // Берем id пост чтобы найти ВСЕ его комментарии
	if err != nil {
		errorW := fmt.Sprint("CommentService.Get: %s", err)
		logger.GetLoggerFromCtx(ctx).Info(ctx, errorW)
		return nil, fmt.Errorf(errorW)
	}
	comments, err := s.repo.GetAllCommentOfPost(ctx, postId) // Находим все комментарии
	if err != nil {
		errorW := fmt.Sprint("CommentService.Get: %s", err)
		logger.GetLoggerFromCtx(ctx).Info(ctx, errorW)
		return nil, fmt.Errorf(errorW)
	}

	var commentbyId []*model.Comment
	for _, j := range comments { // Перебираем все комментарии
		if j.ID == int32(CommentID) { // в качестве корня дерева берем наш целевой комментарий

			commentbyId = append(commentbyId, j)
			var queue []*model.Comment
			queue = FindKids(int(j.ID), comments) // Ищем всех дочерних комментов

			for len(queue) != 0 { // DFS для того чтобы сохранить иерархию коментариев ( BFS не подойдет например )
				tar := queue[len(queue)-1]
				queue = queue[:len(queue)-1]
				commentbyId = append(commentbyId, tar)

				parentTar := FindKids(int(tar.ID), comments)

				queue = append(queue, parentTar...)
			}

		}

	}

	err = graph.CheckLimit(commentbyId, &limit, &offset) //Пигинация
	if err != nil {
		errorW := fmt.Sprint("CommentService.Get: %s", err)
		logger.GetLoggerFromCtx(ctx).Info(ctx, errorW)
		return nil, fmt.Errorf(errorW)
	}
	return commentbyId[offset : offset+limit], nil
}

func (s Service) GetAllPost(ctx context.Context, PostId int, limit int, offset int) ([]*model.Comment, error) {
	//Поиск всех комментариев поста
	// Алгоритм не отличается ничем, за исключением одной строчки
	comments, err := s.repo.GetAllCommentOfPost(ctx, PostId)
	if err != nil {
		errorW := fmt.Sprint("GetAllPost: CommentService.Get: %s", err)
		logger.GetLoggerFromCtx(ctx).Info(ctx, errorW)
		return nil, fmt.Errorf(errorW)
	}

	var preorityComments []*model.Comment
	for _, j := range comments {
		if !CheckSlice(j, preorityComments) && (j.ParentIDComment < 1) { // Здесь в качестве корня берем ВСЕ комметарии у которых нет родителя
			// А тк у нас в списке комментарии комменты только данного поста, у нас получается несколько деревьев, которые выводим
			preorityComments = append(preorityComments, j)
			var queue []*model.Comment
			queue = FindKids(int(j.ID), comments)

			for len(queue) != 0 { //DFS
				tar := queue[len(queue)-1]
				queue = queue[:len(queue)-1]
				preorityComments = append(preorityComments, tar)

				parentTar := FindKids(int(tar.ID), comments)

				queue = append(queue, parentTar...)
			}

		}

	}

	err = graph.CheckLimit(preorityComments, &limit, &offset) //Пигинация
	if err != nil {
		errorW := fmt.Sprint("CommentService.GetAllPost: %s", err)
		logger.GetLoggerFromCtx(ctx).Info(ctx, errorW)
		return nil, fmt.Errorf(errorW)
	}

	return preorityComments[offset : offset+limit], nil
}

func CheckSlice(commentTarget *model.Comment, comments []*model.Comment) bool { // Функция для проверки наличия коментария в списке
	for _, j := range comments {
		if j != nil && commentTarget == j {
			return true
		}
	}
	return false
}
func FindKids(IdParent int, comments []*model.Comment) []*model.Comment { // Поиск из списка всех дочерних ( прямых потомков ) данного коментария
	var res []*model.Comment
	for _, j := range comments {
		if j != nil && (IdParent == int(j.ParentIDComment)) {
			res = append(res, j)
		}
	}
	return res
}
