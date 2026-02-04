package main

import (
	"context"

	router "qraphQL_posts/api"
	"qraphQL_posts/pkg/logger"

	user_repository "qraphQL_posts/internal/User/repository"
	user_service "qraphQL_posts/internal/User/service"

	post_repository "qraphQL_posts/internal/Post/repository"
	post_service "qraphQL_posts/internal/Post/service"

	comment_repository "qraphQL_posts/internal/Comment/repository"
	comment_service "qraphQL_posts/internal/Comment/service"

	localstorage "qraphQL_posts/pkg/localStorage"
)

func main() {
	ctx := context.Background()

	ctx, _ = logger.NewLogger(ctx)
	logger.GetLoggerFromCtx(ctx).Info(ctx, "2")

	localstorage := localstorage.NewLocalStorage()

	UserRepo := user_repository.New(&localstorage)
	UserService := user_service.New(ctx, UserRepo)

	PostRepo := post_repository.New(&localstorage)
	PostService := post_service.New(ctx, PostRepo)

	CommentRepo := comment_repository.New(&localstorage)
	CommentService := comment_service.New(ctx, CommentRepo, PostService)

	router.NewRouter(ctx, UserService, PostService, CommentService)

}
