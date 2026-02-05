package comment_psg_repository

import (
	"context"
	"fmt"
	"qraphQL_posts/api/graph/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

/*
CREATE TABLE IF NOT EXISTS comments (

		id SERIAL PRIMARY KEY,
		post_id INT NOT NULL,
		parent_id INT,
		author_id INT NOT NULL,
		content VARCHAR(2000) NOT NULL,
	);
*/
const (
	queryCreateComment = `
        INSERT INTO comments (post_id, parent_id, author_id, content) 
        VALUES ($1, $2, $3, $4)
        RETURNING id
    `

	queryGetComment = `
        SELECT 
            c.id,
            c.author_id,
            c.content,
            c.parent_id,
            c.post_id,
            c.created_at
        FROM comments c
        WHERE c.post_id = (
            SELECT post_id 
            FROM comments 
            WHERE id = $1
        )
        ORDER BY c.created_at
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

func (r Repository) Create(ctx context.Context, Comment *model.Comment) (*model.Comment, error) {
	var id string
	err := r.pgDB.QueryRow(ctx, queryCreateComment,
		Comment.PostID,
		Comment.ParentIDComment,
		"-67",
		Comment.Content,
	).Scan(&id)

	if err != nil {
		return nil, fmt.Errorf("failed to create comment: %w", err)
	}
	Comment.ID = id
	return Comment, nil

}
func (r Repository) Update(ctx context.Context, Comment *model.Comment) (bool, error) {
	return false, nil
}
func (r Repository) Get(ctx context.Context, CommentId int) (*model.Comment, error) {
	/*
		if CommentId >= len(r.localstorage.Comments) {
			return nil, fmt.Errorf("User does not exist")
		}
			return &r.localstorage.Comments[CommentId], nil
	*/

	rows, err := r.pgDB.Query(ctx, queryGetComment, CommentId)
	if err != nil {
		return nil, fmt.Errorf("failed to get comments hierarchy: %w", err)
	}
	defer rows.Close()

	comments := make(map[int]*model.Comment)
	for rows.Next() {
		var comment *model.Comment
		var temp string
		err := rows.Scan(
			&comment.ID,
			&comment.ParentIDComment,
			&temp,
			&comment.Content,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan comment: %w", err)
		}
		comments[comment.ID] = comment
		comments = append(comments, comment)
	}

	return comments, nil
}
func (r Repository) Delete(ctx context.Context, Post int) (bool, error) {
	return false, nil
}
func (r Repository) GetAllPost(ctx context.Context, PostId int) ([]*model.Comment, error) {
	/*	if PostId >= len(r.localstorage.Posts) || PostId < 0 {
			return nil, fmt.Errorf("User does not exist")
		}
		fmt.Println(r.localstorage.Posts[PostId].Comments)
		return r.localstorage.Posts[PostId].Comments, nil
	*/
	return nil, nil
}
