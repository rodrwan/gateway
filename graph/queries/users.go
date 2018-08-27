package queries

import (
	"errors"

	"github.com/graphql-go/graphql"
	"github.com/rodrwan/gateway"
	"github.com/rodrwan/gateway/graph"
	"github.com/rodrwan/gateway/graph/types"
)

// User fill graphql Field with data from postgres service.
func User(ctx *graph.Context) *graphql.Field {
	return &graphql.Field{
		Type:        types.User,
		Description: "Get user by email",
		Args: graphql.FieldConfigArgument{
			"email": &graphql.ArgumentConfig{
				Type:        graphql.String,
				Description: "return user information by email",
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			email, ok := params.Args["email"].(string)
			if !ok {
				return nil, errors.New("Invalid params")
			}
			// user UserService
			opts := gateway.SetUserQueryOptions(&gateway.UserQueryOptions{
				Email: email,
			})
			u, err := ctx.UserService.Get(opts)
			if err != nil {
				return nil, err
			}

			return u, nil
		},
	}
}

// Users fill graphql Field with data from postgres service.
func Users(ctx *graph.Context) *graphql.Field {
	return &graphql.Field{
		Type:        graphql.NewList(types.User),
		Description: "Get users",
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			uu, err := ctx.UserService.All()
			if err != nil {
				return nil, err
			}

			return uu, nil
		},
	}
}

// UserQuery expose UserQuery
func UserQuery(ctx *graph.Context) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "UserQuery",
		Fields: graphql.Fields{
			"user":  User(ctx),
			"users": Users(ctx),
		},
	})
}
