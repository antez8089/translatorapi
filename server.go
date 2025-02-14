package main

import (
	"log"
	"net/http"
	"translatorapi/database"
	"translatorapi/graph"
	generated "translatorapi/graph/generated"


	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

func main() {
	// Initialize database connection
	db, err  := database.InitDB()

	if err != nil {
		log.Fatalf("Nie udało się połączyć z bazą danych: %v", err)
	}


	// Set up the GraphQL handler with generated executable schema
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{
		Resolvers: &graph.Resolver{DB: db}, // Make sure Resolver is correctly implemented
	}))

	// Serve the GraphQL playground at root
	http.Handle("/", playground.Handler("GraphQL Playground", "/query"))
	// Handle queries at /query
	http.Handle("/query", srv)

	// Start the server
	log.Println("Server running on http://localhost:8080/")
	log.Fatal(http.ListenAndServe(":8080", nil))
}