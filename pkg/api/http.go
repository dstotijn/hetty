package api

import (
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/mux"
)

func HTTPHandler(resolver *Resolver, gqlEndpoint string) http.Handler {
	router := mux.NewRouter().SkipClean(true)
	router.Methods("POST").Handler(
		handler.NewDefaultServer(NewExecutableSchema(Config{
			Resolvers: resolver,
		})),
	)
	router.Methods("GET").Handler(playground.Handler("GraphQL Playground", gqlEndpoint))

	return router
}
