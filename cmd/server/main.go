package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/graphql-go/graphql"
	"github.com/rodrwan/gateway/graph"
	"github.com/rodrwan/gateway/graph/queries"
	"github.com/rodrwan/gateway/postgres"
)

func main() {
	postgresDSN := flag.String(
		"postgres-dsn", "postgres://localhost:5432/graph_test?sslmode=disable", "Postgres DSN")

	flag.Parse()
	db, err := postgres.NewDatastore(*postgresDSN)
	check(err)

	us := &postgres.UserService{
		Store: db,
	}

	ctx := &graph.Context{
		UserService: us,
	}
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: queries.Users(ctx),
	})
	check(err)

	http.HandleFunc("/users", disableCors(getUser(schema)))

	log.Println("Now server is running on port 3000")
	http.ListenAndServe(":3000", nil)
}

// ContentTypeGraphQL graphql content type
const ContentTypeGraphQL = "application/graphql"

func getUser(schema graphql.Schema) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		contentTypeStr := r.Header.Get("Content-Type")
		contentTypeTokens := strings.Split(contentTypeStr, ";")
		contentType := contentTypeTokens[0]

		var result *graphql.Result
		switch contentType {
		case ContentTypeGraphQL:
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Could not read body", http.StatusInternalServerError)
				return
			}

			result = executeQuery(string(body), schema)
		default:
			http.Error(w, "Missing application/graphql", http.StatusBadRequest)
		}

		w.Header().Set("Accept-Encoding", "gzip")
		json.NewEncoder(w).Encode(result)
	}
}

func disableCors(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, Content-Length, Accept-Encoding")
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Max-Age", "86400")
			w.WriteHeader(http.StatusOK)
			return
		}
		h.ServeHTTP(w, r)
	}
}

func executeQuery(query string, schema graphql.Schema) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	if len(result.Errors) > 0 {
		fmt.Printf("errors: %v", result.Errors)
	}
	return result
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
