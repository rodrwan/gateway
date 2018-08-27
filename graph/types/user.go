package types

import "github.com/graphql-go/graphql"

var User = graphql.NewObject(graphql.ObjectConfig{
	Name: "User",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.String,
		},
		"email": &graphql.Field{
			Type: graphql.String,
		},
		"first_name": &graphql.Field{
			Type: graphql.String,
		},
		"last_name": &graphql.Field{
			Type: graphql.String,
		},
		"phone": &graphql.Field{
			Type: graphql.String,
		},
		"birthdate": &graphql.Field{
			Type: graphql.DateTime,
		},
		"address": &graphql.Field{
			Type: Address,
		},
	},
})
