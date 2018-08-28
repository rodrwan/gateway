package users

import (
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/rodrwan/gateway"
	"github.com/rodrwan/gateway/postgres"
	"github.com/rodrwan/gateway/postgres/addresses"
)

// Service implements gateway.UserService using postgres as a backend.
type Service struct {
	Store postgres.SQLExecutor
}

// Get returns the user associated with the given ID.
func (us *Service) Get(opts ...gateway.UserQueryOption) (*gateway.User, error) {
	usstore := UserStore{us.Store}
	astore := addresses.Store{Store: us.Store}

	u, err := usstore.Get(opts...)
	if err != nil {
		return nil, err
	}

	a, err := astore.Find(u.ID)
	if err != nil {
		return nil, err
	}
	u.Address = a

	return u, nil
}

// Select ...
func (us *Service) Select(opts ...gateway.UserQueryOption) ([]*gateway.User, error) {
	usstore := UserStore{us.Store}

	return usstore.Select(opts...)
}

// UsersWithAddress ...
func (us *Service) UsersWithAddress() ([]*gateway.User, error) {
	usstore := UserStore{us.Store}

	return usstore.UsersWithAddress()
}

// All ...
func (us *Service) All() ([]*gateway.User, error) {
	usstore := UserStore{us.Store}
	uu, err := usstore.All()
	if err != nil {
		return nil, err
	}

	astore := addresses.Store{Store: us.Store}
	for _, u := range uu {
		a, err := astore.Find(u.ID)
		if err != nil {
			return nil, err
		}

		u.Address = a
	}
	return uu, nil
}

// Create creates the given user.
func (us *Service) Create(u *gateway.User) error {
	return us.transact(func(store postgres.SQLExecutor) error {
		ustore := UserStore{store}
		u.Phone = strings.Replace(u.Phone, " ", "", -1)

		if err := ustore.Create(u); err != nil {
			return err
		}

		if u.Address == nil {
			return gateway.ErrUserInvalidAddress
		}

		astore := addresses.Store{Store: store}
		u.Address.AddressLine = formatAddress(u.Address.AddressLine)
		u.Address.UserID = u.ID
		return astore.Create(u.Address)
	})
}

// Update updates the fields of the given user. Ignores password and hashed passwords fields.
func (us *Service) Update(u *gateway.User) error {
	return us.transact(func(store postgres.SQLExecutor) error {
		ustore := UserStore{store}

		if err := ustore.Update(u); err != nil {
			return err
		}

		if u.Address == nil {
			return gateway.ErrUserInvalidAddress
		}

		astore := addresses.Store{Store: store}
		u.Address.AddressLine = formatAddress(u.Address.AddressLine)
		u.Address.UserID = u.ID
		return astore.Update(u.Address)
	})
}

// Delete ...
func (us *Service) Delete(u *gateway.User) error {
	return us.transact(func(store postgres.SQLExecutor) error {
		ustore := UserStore{store}
		return ustore.Delete(u)
	})
}

// runs fn in a transactions and do a rollback in case of any error.
func (us *Service) transact(fn func(store postgres.SQLExecutor) error) error {
	var tx *sqlx.Tx
	var err error

	switch store := us.Store.(type) {
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
