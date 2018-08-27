package postgres

import (
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/rodrwan/gateway"
)

type CardStore struct {
	store SQLExecutor
}

// // Get search a user by id.
// func (us *UserStore) Get(opts ...gateway.UserQueryOption) (*gateway.User, error) {
// 	queryOpts := new(gateway.UserQueryOptions)
// 	for _, opt := range opts {
// 		opt(queryOpts)
// 	}

// 	q := squirrel.Select("*").From("users")

// 	if queryOpts.ID != "" {
// 		q = q.Where("id = ?", queryOpts.ID)
// 	}

// 	if queryOpts.Email != "" {
// 		q = q.Where("email = ?", queryOpts.Email)
// 	}

// 	if queryOpts.Token != "" {
// 		q = q.Where("reset_password_token = ?", queryOpts.Token)
// 	}

// 	if queryOpts.DNI != "" {
// 		q = q.Where("dni = ?", queryOpts.DNI)
// 	}

// 	sql, args, err := q.PlaceholderFormat(squirrel.Dollar).ToSql()
// 	if err != nil {
// 		return nil, userError(err)
// 	}

// 	var u gateway.User
// 	row := us.store.QueryRowx(sql, args...)
// 	if err := row.StructScan(&u); err != nil {
// 		fmt.Println(err)
// 		return nil, userError(err)
// 	}

// 	return &u, nil
// }

// Select ...
func (cs *CardStore) Select(opts ...gateway.CardQueryOption) ([]*gateway.Card, error) {
	var opt gateway.CardQueryOptions
	for _, fn := range opts {
		fn(&opt)
	}

	q := squirrel.Select("id, user_id, product_id, pan, ref_id, ref_email, ref_user_id").From("cards").Where("deleted_at is null")

	if opt.UserID != "" {
		q = q.Where("user_id = ?", opt.UserID)
	}

	sql, args, err := q.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := cs.store.Queryx(sql, args...)
	if err != nil {
		return nil, err
	}

	cc := []*gateway.Card{}
	for rows.Next() {
		var c gateway.Card
		if err := rows.StructScan(&c); err != nil {
			return nil, userError(err)
		}
		cc = append(cc, &c)
	}

	return cc, nil
}

func (cs *CardStore) CardDeposits(id string) ([]*gateway.CardDeposit, error) {
	q := squirrel.
		Select("id, amount, payment_id, created_at, status, created_at, fee, total, usd").
		From("card_deposits").
		Where("deleted_at is null")

	q = q.Where("card_id = ?", id)

	sql, args, err := q.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := cs.store.Queryx(sql, args...)
	if err != nil {
		return nil, err
	}

	cc := []*gateway.CardDeposit{}
	for rows.Next() {
		var c gateway.CardDeposit
		if err := rows.StructScan(&c); err != nil {
			return nil, userError(err)
		}
		cc = append(cc, &c)
	}

	return cc, nil
}

type CardService struct {
	Store SQLExecutor
}

// Select ...
func (cs *CardService) Select(opts ...gateway.CardQueryOption) ([]*gateway.Card, error) {
	csstore := CardStore{cs.Store}

	return csstore.Select(opts...)
}

func (cs *CardService) CardDeposits(id string) ([]*gateway.CardDeposit, error) {
	csstore := CardStore{cs.Store}

	return csstore.CardDeposits(id)
}

// runs fn in a transactions and do a rollback in case of any error.
func (cs *CardService) transact(fn func(store SQLExecutor) error) error {
	var tx *sqlx.Tx
	var err error

	switch store := cs.Store.(type) {
	case *sqlx.DB:
		tx, err = store.Beginx()
		if err != nil {
			return err
		}
	case *sqlx.Tx:
		tx = store
	}

	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
