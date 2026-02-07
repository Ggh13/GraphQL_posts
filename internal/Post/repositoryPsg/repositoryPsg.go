package post_psg_repository

import (
	"context"
	"fmt"
	"qraphQL_posts/api/graph/model"
	"qraphQL_posts/pkg/logger"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	queryCreatePost = `
        INSERT INTO posts (content, author_id, commentable) 
        VALUES ($1, $2, $3)
        RETURNING id
    `
	queryCheckUser = `
		 SELECT 
			id
		FROM users
		WHERE id = $1
	`
	queryGetAuthorId = `
		 SELECT 
			author_id
		FROM posts
		WHERE id = $1
	`
	queryGetPost = `
        SELECT 
			p.id,
			p.content,
			p.commentable,
			u.id,
			u.name ,
			u.surname
		FROM posts p
		LEFT JOIN users u ON p.author_id = u.id
		WHERE p.id = $1
		ORDER BY p.id
    `

	queryUpdatePos = `
		UPDATE posts 
		SET commentable = $1 
		WHERE id = $2
	`
	queryGetAllPosts = `
    SELECT 
        p.id,
        p.content,
        p.commentable,
        u.id,
        u.name,
        u.surname
    FROM posts p
    LEFT JOIN users u ON p.author_id = u.id
    ORDER BY p.id
`
)

type Repository struct {
	pgDB *pgxpool.Pool
}

func New(ctx context.Context, pgDB *pgxpool.Pool) *Repository {
	return &Repository{pgDB: pgDB}
}
func (r Repository) Create(ctx context.Context, Post *model.Post) (*model.Post, error) {
	var id_user int
	id_user = int(Post.User.ID) //Проверяем существование user которые создает пост
	err := r.pgDB.QueryRow(ctx, queryCheckUser,
		&id_user,
	).Scan(&id_user)
	if err != nil {
		errorW := fmt.Sprintf("PostRepository.Create: failed to create post, user with this id does not exist: %v", err)
		logger.GetLoggerFromCtx(ctx).Info(ctx, errorW)
		return nil, fmt.Errorf(errorW)
	}

	var id string
	err = r.pgDB.QueryRow(ctx, queryCreatePost, //Создаем пост
		Post.Content,
		Post.User.ID,
		Post.Commentable,
	).Scan(&id)

	if err != nil {
		errorW := fmt.Sprintf("PostRepository.Create: failed to create post: %v", err)
		logger.GetLoggerFromCtx(ctx).Info(ctx, errorW)
		return nil, fmt.Errorf(errorW)
	}
	idInt, err := strconv.Atoi(id)
	if err != nil {
		errorW := fmt.Sprintf("PostRepository.Create: failed to create post: %v", err)
		logger.GetLoggerFromCtx(ctx).Info(ctx, errorW)
		return nil, fmt.Errorf(errorW)
	}
	Post.ID = int32(idInt)
	return Post, nil
}

func (r Repository) Update(ctx context.Context, Post *model.Post) (bool, error) {

	var authorId int
	err := r.pgDB.QueryRow(ctx, queryGetAuthorId, int(Post.ID)).Scan( // Получаем автора поста
		&authorId,
	)

	if err != nil {
		errorW := fmt.Sprintf("PostRepository.Update: Failed find post author by ID %d %v", Post.ID, err)
		logger.GetLoggerFromCtx(ctx).Info(ctx, errorW)
		return false, fmt.Errorf(errorW)
	}

	if authorId != int(Post.User.ID) { // проверяем является ли user обновляющий данные автором
		errorW := fmt.Sprintf("PostRepository.Update: You are not author of this post. You can not update it")
		logger.GetLoggerFromCtx(ctx).Info(ctx, errorW)
		return false, fmt.Errorf(errorW)
	}
	_, err = r.pgDB.Exec(context.Background(), queryUpdatePos, // Обновляем ТОЛЬКО право комментирования (по ТЗ)
		Post.Commentable, Post.ID,
	)

	if err != nil {
		errorW := fmt.Sprintf("PostRepository.Update: Cant update post with id %v %v", Post.ID, err)
		logger.GetLoggerFromCtx(ctx).Info(ctx, errorW)
		return false, fmt.Errorf(errorW)
	}

	return true, nil
}
func (r Repository) Get(ctx context.Context, PostId int) (*model.Post, error) {
	var post model.Post
	var user model.User
	err := r.pgDB.QueryRow(ctx, queryGetPost, PostId).Scan(
		&post.ID,
		&post.Content,
		&post.Commentable,
		&user.ID,
		&user.Name,
		&user.Surname,
	)
	post.User = &user
	if err != nil {
		errorW := fmt.Sprintf("PostRepository.Get: Failed to get post by ID %s %v", PostId, err)
		logger.GetLoggerFromCtx(ctx).Info(ctx, errorW)
		return nil, fmt.Errorf(errorW)
	}

	return &post, nil
}
func (r Repository) Delete(ctx context.Context, Post int) (bool, error) {
	return false, nil
}

func (r Repository) GetAllPosts(ctx context.Context) ([]*model.Post, error) {
	rows, err := r.pgDB.Query(ctx, queryGetAllPosts)
	if err != nil {
		errorW := fmt.Sprintf("PostRepository.GetAllPosts: %v", err)
		logger.GetLoggerFromCtx(ctx).Info(ctx, errorW)
		return nil, fmt.Errorf(errorW)
	}
	defer rows.Close()

	if rows == nil {
		return nil, nil
	}

	var posts []*model.Post
	for rows.Next() { // Запрос не позволяющий появится проблеме N+1 за счет JOIN с таблицей user
		var post model.Post
		var user model.User
		err := rows.Scan(
			&post.ID,
			&post.Content,
			&post.Commentable,
			&user.ID,
			&user.Name,
			&user.Surname,
		)
		if err != nil {
			errorW := fmt.Sprintf("PostRepository.GetAllPosts: failed to scan posts: %v", err)
			logger.GetLoggerFromCtx(ctx).Info(ctx, errorW)
			return nil, fmt.Errorf(errorW)
		}
		post.User = &user
		posts = append(posts, &post)

	}

	return posts, nil
}
