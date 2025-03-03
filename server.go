package main

import (
	"log"
	"net/http"
	"translatorapi/database"
	"translatorapi/graph"

	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/vektah/gqlparser/v2/ast"

	generated "translatorapi/graph/generated"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

func main() {
	// Initialize database connection
	db, err := database.InitDB()

	if err != nil {
		log.Fatalf("Nie udało się połączyć z bazą danych: %v", err)
	}

	// // Set up the GraphQL handler with generated executable schema
	// srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{
	// 	Resolvers: &graph.Resolver{DB: db}, // Make sure Resolver is correctly implemented
	// }))

	srv := handler.New(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{DB: db}}))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	// Serve the GraphQL playground at root
	http.Handle("/", playground.Handler("GraphQL Playground", "/query"))
	// Handle queries at /query
	http.Handle("/query", srv)

	// Start the server
	log.Println("Server running on http://localhost:8080/")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
