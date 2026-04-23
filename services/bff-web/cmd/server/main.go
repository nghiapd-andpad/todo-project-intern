package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/di"
	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/config"
	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/handler/dataloader"
	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/handler/directive"
	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/handler/graph"
	authmw "github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/handler/http/middleware"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	app, cleanup, err := di.InitializeApp(cfg)
	if err != nil {
		log.Fatalf("failed to initialize app: %v", err)
	}
	defer cleanup()

	gqlConfig := graph.Config{Resolvers: app.Resolver}
	gqlConfig.Directives.Auth = directive.Auth

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(gqlConfig))
	srv.SetErrorPresenter(graph.ErrorPresenter)

	mux := http.NewServeMux()
	mux.Handle("/", playground.Handler("GraphQL Playground", "/graphql"))
	mux.Handle("/graphql", authmw.AuthMiddleware(app.JwtManager)(
		dataloader.Middleware(app.Resolver.GetUserGetter())(srv),
	))

	addr := ":" + cfg.ServerPort
	fmt.Printf("BFF GraphQL server listening on %s\n", addr)
	fmt.Printf("Playground: http://localhost%s\n", addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
