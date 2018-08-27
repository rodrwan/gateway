package types

import (
	"github.com/graphql-go/graphql"
	"github.com/rodrwan/gateway"
)

var Address = graphql.NewObject(graphql.ObjectConfig{
	Name: "Address",
	Fields: graphql.Fields{
		"user_id": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				a := p.Source.(*gateway.Address)
				return a.City, nil
			},
		},
		"city": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				a := p.Source.(*gateway.Address)
				return a.City, nil
			},
		},
		"address_line": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				a := p.Source.(*gateway.Address)
				return a.AddressLine, nil
			},
		},
		"locality": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				a := p.Source.(*gateway.Address)
				return a.Locality, nil
			},
		},
		"administrative_area_level_1": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				a := p.Source.(*gateway.Address)
				return a.AdministrativeAreaLevel1, nil
			},
		},
		"country": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				a := p.Source.(*gateway.Address)
				return a.Country, nil
			},
		},
		"postal_code": &graphql.Field{
			Type: graphql.Int,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				a := p.Source.(*gateway.Address)
				return a.PostalCode, nil
			},
		},
	},
})
