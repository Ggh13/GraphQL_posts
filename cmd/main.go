package main

import (
	"context"
	router "qraphQL_posts/api"
	"qraphQL_posts/pkg/logger"
)

func main() {
	ctx := context.Background()

	ctx, _ = logger.NewLogger(ctx)

	router.NewRouter(ctx)

}
