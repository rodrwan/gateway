package queries

// // Field User.cards
// types.User.AddFieldConfig("cards", UserCards(ctx))
// // Field Card.deposits
// types.Card.AddFieldConfig("deposits", CardDeposits(ctx))

// // UserCards field that expose cards of a user by id
// func UserCards(ctx *graph.Context) *graphql.Field {
// 	return &graphql.Field{
// 		Type:        graphql.NewList(types.Card),
// 		Description: "Get cards by user id",
// 		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
// 			u := p.Source.(*gateway.User)
// 			opts := gateway.SetCardQueryOptions(&gateway.CardQueryOptions{
// 				UserID: u.ID,
// 			})
// 			cc, err := ctx.CardService.Select(opts)
// 			if err != nil {
// 				return nil, err
// 			}

// 			return cc, nil
// 		},
// 	}
// }
