package types

import (
	"github.com/graphql-go/graphql"
)

var CardDeposit = graphql.NewObject(graphql.ObjectConfig{
	Name: "CardDeposit",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.String,
		},
		"amount": &graphql.Field{
			Type: graphql.Int,
		},
		"payment_id": &graphql.Field{
			Type: graphql.String,
		},
		"status": &graphql.Field{
			Type: graphql.String,
		},
		"created_at": &graphql.Field{
			Type: graphql.DateTime,
		},
		"fee": &graphql.Field{
			Type: graphql.Int,
		},
		"total": &graphql.Field{
			Type: graphql.Int,
		},
		"dollar": &graphql.Field{
			Type: graphql.Int,
		},
	},
})
