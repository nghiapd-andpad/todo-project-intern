package main

import (
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/di"
	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/handler/directive"
	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/handler/graph"
	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/handler/http/middleware"
)

func main() {
	resolver, cleanup, err := di.InitializeResolver()
	if err != nil {
		log.Fatalf("failed to initialize resolver: %v", err)
	}
	defer cleanup()

	// Config GraphQL with Directives (@auth, @hasRoles)
	graphConfig := graph.Config{Resolvers: resolver}
	graphConfig.Directives.Auth = directive.Auth
	graphConfig.Directives.HasRoles = directive.HasRoles

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graphConfig))

	srv.SetErrorPresenter(graph.PresentError)

	authUC := resolver.AuthUseCase

	http.Handle("/", playground.Handler("GraphQL playground", "/graphql"))

	http.Handle("/graphql", middleware.AuthMiddleware(authUC)(srv))

	log.Printf("BFF Web Server is running on port :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
