package main

import (
	"context"
	"fmt"
	router "qraphQL_posts/api"
	"qraphQL_posts/pkg/logger"

	user_repository "qraphQL_posts/internal/User/repository"
	user_service "qraphQL_posts/internal/User/service"
	localstorage "qraphQL_posts/pkg/localStorage"
)

func main() {
	ctx := context.Background()

	ctx, _ = logger.NewLogger(ctx)
	logger.GetLoggerFromCtx(ctx).Info(ctx, "2")
	fmt.Println(logger.GetLoggerFromCtx(ctx))
	localstorage := localstorage.NewLocalStorage()
	fmt.Println(logger.GetLoggerFromCtx(ctx))
	UserRepo := user_repository.New(&localstorage)
	fmt.Println(logger.GetLoggerFromCtx(ctx))
	UserService := user_service.New(ctx, UserRepo)
	fmt.Println(logger.GetLoggerFromCtx(ctx))
	router.NewRouter(ctx, UserService, nil, nil)

}
