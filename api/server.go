package router

import (
	"context"
	"fmt"
	"net/http"
	"qraphQL_posts/api/graph"
	"qraphQL_posts/pkg/logger"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

func NewRouter(ctx context.Context, resolver *graph.Resolver, port string, srv *handler.Server) {

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		logger.GetLoggerFromCtx(ctx).Fatal(ctx,
			fmt.Sprintf("failed to start server: %v", err))
	}

}
