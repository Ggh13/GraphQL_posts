package comment_psg_repository

import (
	"context"
	"fmt"
	"qraphQL_posts/api/graph/model"
	"qraphQL_posts/pkg/logger"
	"strconv"

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

	queryGetCommentByParent = `
		SELECT 
			id
		FROM comments 
		WHERE id = $1
		ORDER BY id
`

	queryGetComment = `
		SELECT 
			c.id,
			c.content,
			c.parent_id,
			c.post_id,
			u.id,
			u.name,
			u.surname,
			u.email
		FROM comments c
		LEFT JOIN users u ON c.author_id = u.id
		WHERE c.post_id = (
			SELECT post_id 
			FROM comments 
			WHERE id = $1
		)
		ORDER BY c.id
`

	queryGetAllCommentOfPost = `
        SELECT 
			c.id,
			c.content,
			c.parent_id,
			c.post_id,
			u.id,
			u.name,
			u.surname
		FROM comments c
		LEFT JOIN users u ON c.author_id = u.id
		WHERE c.post_id = $1
		ORDER BY c.id
    `
	queryGetIdCommbyPostId = `
        SELECT 
            post_id
        FROM comments 
        WHERE id = $1
		ORDER BY id
    `

	queryCheckUser = `
		 SELECT 
			id
		FROM users
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

func (r Repository) Create(ctx context.Context, Comment *model.Comment) (*model.Comment, error) {

	var id_user int
	id_user = int(Comment.User.ID)
	err := r.pgDB.QueryRow(ctx, queryCheckUser,
		&id_user,
	).Scan(&id_user)
	if err != nil {
		return nil, fmt.Errorf("failed to create post, user with this id does not exist: %w", err)
	}

	if Comment.ParentIDComment >= 1 {
		var id_parent string
		err := r.pgDB.QueryRow(ctx, queryGetCommentByParent,
			Comment.ParentIDComment,
		).Scan(&id_parent)

		if err != nil {
			return nil, fmt.Errorf("failed to create comment. Parent comment must exist: %w", err)
		}
	}

	var id string
	err = r.pgDB.QueryRow(ctx, queryCreateComment,
		Comment.PostID,
		Comment.ParentIDComment,
		Comment.User.ID,
		Comment.Content,
	).Scan(&id)

	if err != nil {
		return nil, fmt.Errorf("failed to create comment: %w", err)
	}
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, fmt.Errorf("Error in Create comment")
	}

	Comment.ID = int32(idInt)

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
		return nil, fmt.Errorf("failed to get comments %w", err)
	}
	defer rows.Close()

	if rows == nil {
		return nil, nil
	}

	// Fisrt stage. Get ALL comments of this post
	comments := make(map[int]*model.Comment)
	for rows.Next() {
		var comment model.Comment
		err := rows.Scan(
			&comment.ID,
			&comment.Content,
			&comment.ParentIDComment,
			&comment.PostID,
			&comment.User.ID,
			&comment.User.Name,
			&comment.User.Surname,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan comment: %w", err)
		}
		comments[int(comment.ID)] = &comment

	}

	// First ( with half ) stage.
	for _, j := range comments {
		if int(j.ParentIDComment) >= 1 {
			comments[int(j.ParentIDComment)].Comments = append(comments[int(j.ParentIDComment)].Comments, j)
		}

	}
	return comments[CommentId], nil

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
				queryGetAllCommentOfPost = `
		        SELECT
					c.id,
					c.content,
					c.parent_id,
					c.post_id,
					u.id as user_id,
					u.name as user_name,
					u.surname as user_surname
				FROM comments c
				LEFT JOIN users u ON c.author_id = u.id
				WHERE c.post_id = $1
				ORDER BY c.id
		    `
	*/

	rows, err := r.pgDB.Query(ctx, queryGetAllCommentOfPost, PostId)
	if err != nil {
		return nil, fmt.Errorf("failed to get comments hierarchy: %w", err)
	}
	defer rows.Close()

	// Fisrt stage. Get ALL comments of this post
	var res []*model.Comment
	for rows.Next() {
		var comment model.Comment
		var user model.User
		err := rows.Scan(
			&comment.ID,
			&comment.Content,
			&comment.ParentIDComment,
			&comment.PostID,
			&user.ID,
			&user.Name,
			&user.Surname,
		)
		comment.User = &user
		logger.GetLoggerFromCtx(ctx).Info(ctx, fmt.Sprintf("There are com %v ", comment.ID))
		if err != nil {
			return nil, fmt.Errorf("failed to scan comment: %w", err)
		}
		res = append(res, &comment)

	}

	logger.GetLoggerFromCtx(ctx).Info(ctx, fmt.Sprintf("There are len  %v ", len(res)))
	return res, nil
}

func (r Repository) GetPostId(ctx context.Context, CommentId int) (int, error) {

	var postId int
	err := r.pgDB.QueryRow(ctx, queryGetIdCommbyPostId, CommentId).Scan(
		&postId,
	)

	if err != nil {
		return -1, fmt.Errorf("Failed to get post by ID comment %s %w", CommentId, err)
	}

	return postId, nil
}
