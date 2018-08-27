package types

import "github.com/graphql-go/graphql"

var Card = graphql.NewObject(graphql.ObjectConfig{
	Name: "Card",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.String,
		},
		"user_id": &graphql.Field{
			Type: graphql.String,
		},
		"product_id": &graphql.Field{
			Type: graphql.String,
		},
		"card_number": &graphql.Field{
			Type: graphql.String,
		},
		"reference_id": &graphql.Field{
			Type: graphql.String,
		},
		"reference_email": &graphql.Field{
			Type: graphql.String,
		},
		"reference_user_id": &graphql.Field{
			Type: graphql.String,
		},
	},
})
