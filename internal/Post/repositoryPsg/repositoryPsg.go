package post_psg_repository

import (
	"context"
	"fmt"
	"qraphQL_posts/api/graph/model"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
)

/*
CREATE TABLE IF NOT EXISTS posts (

		id SERIAL PRIMARY KEY,
		content TEXT NOT NULL,
		author_id INT NOT NULL,
		commentable BOOLEAN DEFAULT TRUE
	);
*/
const (
	queryCreatePost = `
        INSERT INTO posts (content, author_id, commentable) 
        VALUES ($1, $2, $3)
        RETURNING id
    `

	queryGetPost = `
        SELECT id, content, author_id, commentable
        FROM posts 
        WHERE id = $1
    `
)

type Repository struct {
	pgDB *pgxpool.Pool
}

/*
	type Repository interface {
		Get(ctx context.Context) (bool, error)
		Create(ctx context.Context, User *model.User) (bool, error)
		Update(ctx context.Context, User *model.User) (bool, error)
		Delete(ctx context.Context, User *model.User) (bool, error)
	}
*/

func New(ctx context.Context, pgDB *pgxpool.Pool) *Repository {
	return &Repository{pgDB: pgDB}
}

func (r Repository) Create(ctx context.Context, Post *model.Post) (*model.Post, error) {
	/*
		Post.ID = int32(len(r.localstorage.Posts))
		r.localstorage.Posts = append(r.localstorage.Posts, *Post)
		return Post, nil
	*/

	var id string
	err := r.pgDB.QueryRow(ctx, queryCreatePost,
		Post.Content,
		Post.UserID,
		Post.Commentable,
	).Scan(&id)

	if err != nil {
		return nil, fmt.Errorf("failed to create post: %w", err)
	}
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, fmt.Errorf("failed to create post: %w", err)
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
	return false, nil
}
func (r Repository) Get(ctx context.Context, PostId int) (*model.Post, error) {
	/*
		if PostId >= len(r.localstorage.Posts) {
			return nil, fmt.Errorf("User does not exist")
		}
		return &r.localstorage.Posts[PostId], nil
	*/
	var post model.Post
	err := r.pgDB.QueryRow(ctx, queryGetPost, PostId).Scan(
		&post.ID,
		&post.Content,
		&post.UserID,
		&post.Commentable,
	)

	if err != nil {
		return nil, fmt.Errorf("Failed to get post by ID %s %w", PostId, err)
	}

	return &post, nil
}
func (r Repository) Delete(ctx context.Context, Post int) (bool, error) {
	return false, nil
}
