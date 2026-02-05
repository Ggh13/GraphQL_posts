package main

import (
	"context"

	router "qraphQL_posts/api"
	"qraphQL_posts/pkg/logger"
	"qraphQL_posts/pkg/postgres"

	user_repository "qraphQL_posts/internal/User/repository"
	user_service "qraphQL_posts/internal/User/service"
	"qraphQL_posts/internal/config"

	post_repository "qraphQL_posts/internal/Post/repository"
	post_service "qraphQL_posts/internal/Post/service"

	comment_repository "qraphQL_posts/internal/Comment/repository"
	comment_service "qraphQL_posts/internal/Comment/service"

	localstorage "qraphQL_posts/pkg/localStorage"

	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()

	ctx, _ = logger.NewLogger(ctx)

	config, err := config.NewConfig(ctx)

	if err != nil {
		logger.GetLoggerFromCtx(ctx).Fatal(ctx, "Failed load config", zap.Error(err))
	}
	logger.GetLoggerFromCtx(ctx).Info(ctx, "Succesfully load config")

	pgDB, err := postgres.NewPostgres(ctx, &config.PostgresCFG)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Fatal(ctx, "Failsed connect to postgres DB", zap.Error(err))
	}
	if err := pgDB.Ping(ctx); err != nil {
		logger.GetLoggerFromCtx(ctx).Fatal(ctx, "Failed ping pgDB", zap.Error(err))
	}
	logger.GetLoggerFromCtx(ctx).Info(ctx, "Succesfully connected to pgDB")

	localstorage := localstorage.NewLocalStorage()

	UserRepo := user_repository.New(&localstorage)
	UserService := user_service.New(ctx, UserRepo)

	PostRepo := post_repository.New(&localstorage)
	PostService := post_service.New(ctx, PostRepo)

	CommentRepo := comment_repository.New(&localstorage)
	CommentService := comment_service.New(ctx, CommentRepo, PostService)

	router.NewRouter(ctx, UserService, PostService, CommentService)

}
