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
	id_user = int(Post.User.ID)
	err := r.pgDB.QueryRow(ctx, queryCheckUser,
		&id_user,
	).Scan(&id_user)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Info(ctx, fmt.Sprint("PostRepository.Create: failed to create post, user with this id does not exist: %w", err))
		return nil, fmt.Errorf("PostRepository.Create: failed to create post, user with this id does not exist: %w", err)
	}

	var id string
	err = r.pgDB.QueryRow(ctx, queryCreatePost,
		Post.Content,
		Post.User.ID,
		Post.Commentable,
	).Scan(&id)

	if err != nil {
		logger.GetLoggerFromCtx(ctx).Info(ctx, fmt.Sprint("PostRepository.Create: failed to create post: %w", err))
		return nil, fmt.Errorf("PostRepository.Create: failed to create post: %w", err)
	}
	idInt, err := strconv.Atoi(id)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Info(ctx, fmt.Sprint("PostRepository.Create: failed to create post: %w", err))
		return nil, fmt.Errorf("PostRepository.Create: failed to create post: %w", err)
	}
	Post.ID = int32(idInt)
	return Post, nil
}

func (r Repository) Update(ctx context.Context, Post *model.Post) (bool, error) {
	/*
		idi := Post.ID
		if idi > int32(len(r.localstorage.Posts)-2) {
			return false, fmt.Errorf("The id post does not exist")
		}
		r.localstorage.Posts[idi].Commentable = *&Post.Commentable
		return true, nil
	*/

	_, err := r.pgDB.Exec(context.Background(), queryUpdatePos,
		Post.Commentable, Post.ID,
	)

	if err != nil {
		logger.GetLoggerFromCtx(ctx).Info(ctx, fmt.Sprint("PostRepository.Update: Cant update post with id %v %w", Post.ID, err))
		return false, fmt.Errorf("PostRepository.Update: Cant update post with id %v %w", Post.ID, err)
	}

	return true, nil
}
func (r Repository) Get(ctx context.Context, PostId int) (*model.Post, error) {
	/*
				if PostId >= len(r.localstorage.Posts) {
					return nil, fmt.Errorf("User does not exist")
				}
				return &r.localstorage.Posts[PostId], nil
				queryGetPost = `
		        SELECT
					p.id,
					p.content,
					p.commentable,
					u.id as id,
					u.name as name,
					u.surname as surname
				FROM posts p
				LEFT JOIN users u ON p.author_id = u.id
				WHERE p.id = $1
		    `
	*/
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
		logger.GetLoggerFromCtx(ctx).Info(ctx, fmt.Sprint("PostRepository.Get: Failed to get post by ID %s %w", PostId, err))
		return nil, fmt.Errorf("PostRepository.Get: Failed to get post by ID %s %w", PostId, err)
	}

	return &post, nil
}
func (r Repository) Delete(ctx context.Context, Post int) (bool, error) {
	return false, nil
}

func (r Repository) GetAllPosts(ctx context.Context) ([]*model.Post, error) {

	/*
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
	*/
	rows, err := r.pgDB.Query(ctx, queryGetAllPosts)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Info(ctx, fmt.Sprint("PostRepository.GetAllPosts: %s", err))
		return nil, fmt.Errorf("PostRepository.GetAllPosts: %s", err)
	}
	defer rows.Close()

	if rows == nil {
		return nil, nil
	}

	// Fisrt stage. Get ALL comments of this post
	var posts []*model.Post
	for rows.Next() {
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
			return nil, fmt.Errorf("failed to scan posts: %w", err)
		}
		post.User = &user
		posts = append(posts, &post)

	}

	return posts, nil
}
