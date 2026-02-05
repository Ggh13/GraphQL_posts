package main

import (
	"context"
	"fmt"

	router "qraphQL_posts/api"
	"qraphQL_posts/pkg/logger"
	"qraphQL_posts/pkg/postgres"

	user_repository "qraphQL_posts/internal/User/repository"
	user_psg_repository "qraphQL_posts/internal/User/repositoryPsg"
	user_service "qraphQL_posts/internal/User/service"
	"qraphQL_posts/internal/config"

	post_repository "qraphQL_posts/internal/Post/repository"
	post_psg_repository "qraphQL_posts/internal/Post/repositoryPsg"
	post_service "qraphQL_posts/internal/Post/service"

	comment_repository "qraphQL_posts/internal/Comment/repository"
	comment_psg_repository "qraphQL_posts/internal/Comment/repositoryPsg"
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

	localstorage := localstorage.NewLocalStorage()

	var UserRepo user_service.Repository
	var PostRepo post_service.Repository
	var CommentRepo comment_service.Repository

	logger.GetLoggerFromCtx(ctx).Info(ctx, fmt.Sprintf("Type of DB was choosen %v", config.TypeDB))

	if config.TypeDB == "PSG" {
		logger.GetLoggerFromCtx(ctx).Info(ctx, "Type of DB was choosen PostgreSQL")

		pgDB, err := postgres.NewPostgres(ctx, &config.PostgresCFG)
		if err != nil {
			logger.GetLoggerFromCtx(ctx).Fatal(ctx, "Failsed connect to postgres DB", zap.Error(err))
		}
		if err := pgDB.Ping(ctx); err != nil {
			logger.GetLoggerFromCtx(ctx).Fatal(ctx, "Failed ping pgDB", zap.Error(err))
		}
		logger.GetLoggerFromCtx(ctx).Info(ctx, "Succesfully connected to pgDB")

		UserRepo = user_psg_repository.New(ctx, pgDB)
		PostRepo = post_psg_repository.New(ctx, pgDB)
		CommentRepo = comment_psg_repository.New(ctx, pgDB)
	} else {
		logger.GetLoggerFromCtx(ctx).Info(ctx, "Type of DB was choosen IN_MEMORY")

		UserRepo = user_repository.New(&localstorage)
		PostRepo = post_repository.New(&localstorage)
		CommentRepo = comment_repository.New(&localstorage)
	}

	UserService := user_service.New(ctx, UserRepo)
	PostService := post_service.New(ctx, PostRepo)
	CommentService := comment_service.New(ctx, CommentRepo, PostService)

	router.NewRouter(ctx, UserService, PostService, CommentService)

}
