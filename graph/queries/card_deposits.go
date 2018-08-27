package queries

// // CardDeposits field that expose card deposits of an specific card
// func CardDeposits(ctx *graph.Context) *graphql.Field {
// 	return &graphql.Field{
// 		Type:        graphql.NewList(types.CardDeposit),
// 		Description: "Get deposits by card id",
// 		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
// 			c := p.Source.(*gateway.Card)
// 			cds, err := ctx.CardService.CardDeposits(c.ID)
// 			if err != nil {
// 				return nil, err
// 			}

// 			return cds, nil
// 		},
// 	}
// }
