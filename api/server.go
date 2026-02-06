package router

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"qraphQL_posts/api/graph"
	"qraphQL_posts/pkg/logger"
	unitTest "qraphQL_posts/test"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/vektah/gqlparser/v2/ast"
)

const defaultPort = "8080"

func NewRouter(ctx context.Context, UserService graph.UserService, PostService graph.PostService, CommentService graph.CommentService) graph.Resolver {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	resolver := &graph.Resolver{}
	fmt.Println("check router")
	fmt.Println(logger.GetLoggerFromCtx(ctx))
	resolver = resolver.NewResolver(ctx, UserService, PostService, CommentService)
	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	unitTest.MainTest(ctx, *resolver)

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
	return *resolver
}
