package postgres

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/rodrwan/gateway"
)

const (
	addrMaxLen = 30
)

// UserStore implements gateway.UserStore interface with postgres as backend
type UserStore struct {
	store SQLExecutor
}

// Get search a user by id.
func (us *UserStore) Get(opts ...gateway.UserQueryOption) (*gateway.User, error) {
	queryOpts := new(gateway.UserQueryOptions)
	for _, opt := range opts {
		opt(queryOpts)
	}

	q := squirrel.Select("*").From("users")

	if queryOpts.ID != "" {
		q = q.Where("id = ?", queryOpts.ID)
	}

	if queryOpts.Email != "" {
		q = q.Where("email = ?", queryOpts.Email)
	}

	if queryOpts.FirstName != "" {
		q = q.Where("first_name = ?", queryOpts.FirstName)
	}

	if queryOpts.LastName != "" {
		q = q.Where("last_name = ?", queryOpts.LastName)
	}

	sql, args, err := q.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return nil, userError(err)
	}

	var u gateway.User
	row := us.store.QueryRowx(sql, args...)
	if err := row.StructScan(&u); err != nil {
		fmt.Println(err)
		return nil, userError(err)
	}

	return &u, nil
}

// Select ...
func (us *UserStore) Select(opts ...gateway.UserQueryOption) ([]*gateway.User, error) {
	var opt gateway.UserQueryOptions
	for _, fn := range opts {
		fn(&opt)
	}

	q := squirrel.Select("*").From("users").Where("deleted_at is null")

	sql, args, err := q.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := us.store.Queryx(sql, args...)
	if err != nil {
		return nil, err
	}

	uu := []*gateway.User{}
	for rows.Next() {
		var u gateway.User
		if err := rows.StructScan(&u); err != nil {
			return nil, userError(err)
		}
		uu = append(uu, &u)
	}

	return uu, nil
}

// All return all active users.
func (us *UserStore) All() ([]*gateway.User, error) {
	rows, err := us.store.Queryx(
		"select * from users",
	)
	if err != nil {
		return nil, userError(err)
	}
	defer rows.Close()

	users := []*gateway.User{}
	for rows.Next() {
		var user gateway.User
		if err := rows.StructScan(&user); err != nil {
			return nil, userError(err)
		}
		users = append(users, &user)
	}

	return users, nil
}

// Create creates a new user
func (us *UserStore) Create(u *gateway.User) error {
	sql, args, err := squirrel.Insert("users").
		Columns(
			"email", "first_name", "last_name", "phone", "birthdate",
		).
		Values(
			u.Email, u.FirstName, u.LastName, u.Phone, u.Birthdate,
		).
		Suffix("returning *").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return userError(err)
	}

	row := us.store.QueryRowx(sql, args...)
	if err := row.StructScan(u); err != nil {
		return userError(err)
	}
	return nil
}

// Update update the given user.
func (us *UserStore) Update(u *gateway.User) error {
	sql, args, err := squirrel.Update("users").
		Set("email", u.Email).
		Set("first_name", u.FirstName).
		Set("last_name", u.LastName).
		Set("phone", u.Phone).
		Set("birthdate", u.Birthdate).
		Where("id = ?", u.ID).
		Suffix("returning *").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return userError(err)
	}

	row := us.store.QueryRowx(sql, args...)
	if err := row.StructScan(u); err != nil {
		return userError(err)
	}

	return nil
}

// Delete mark a user as deleted (logical delete).
func (us *UserStore) Delete(u *gateway.User) error {
	row := us.store.QueryRowx(
		"update users set deleted_at = $1 where id = $2 returning *",
		time.Now(), u.ID,
	)

	if err := row.StructScan(u); err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return err
	}

	return nil
}

// UserService implements gateway.UserService using postgres as a backend.
type UserService struct {
	Store SQLExecutor
}

// Get returns the user associated with the given ID.
func (us *UserService) Get(opts ...gateway.UserQueryOption) (*gateway.User, error) {
	usstore := UserStore{us.Store}
	astore := AddressStore{us.Store}

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
func (us *UserService) Select(opts ...gateway.UserQueryOption) ([]*gateway.User, error) {
	usstore := UserStore{us.Store}

	return usstore.Select(opts...)
}

// All ...
func (us *UserService) All() ([]*gateway.User, error) {
	usstore := UserStore{us.Store}
	uu, err := usstore.All()
	if err != nil {
		return nil, err
	}

	astore := AddressStore{us.Store}
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
func (us *UserService) Create(u *gateway.User) error {
	return us.transact(func(store SQLExecutor) error {
		ustore := UserStore{store}
		u.Phone = strings.Replace(u.Phone, " ", "", -1)

		if err := ustore.Create(u); err != nil {
			return err
		}

		if u.Address == nil {
			return gateway.ErrUserInvalidAddress
		}

		astore := AddressStore{store}
		u.Address.AddressLine = formatAddress(u.Address.AddressLine)
		u.Address.UserID = u.ID
		return astore.Create(u.Address)
	})
}

// Update updates the fields of the given user. Ignores password and hashed passwords fields.
func (us *UserService) Update(u *gateway.User) error {
	return us.transact(func(store SQLExecutor) error {
		ustore := UserStore{store}
		astore := AddressStore{store}

		if err := ustore.Update(u); err != nil {
			return err
		}

		if u.Address == nil {
			return gateway.ErrUserInvalidAddress
		}

		u.Address.AddressLine = formatAddress(u.Address.AddressLine)
		u.Address.UserID = u.ID
		return astore.Update(u.Address)
	})
}

// Delete ...
func (us *UserService) Delete(u *gateway.User) error {
	return us.transact(func(store SQLExecutor) error {
		ustore := UserStore{store}
		return ustore.Delete(u)
	})
}

// runs fn in a transactions and do a rollback in case of any error.
func (us *UserService) transact(fn func(store SQLExecutor) error) error {
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

func userError(err error) error {
	if err == sql.ErrNoRows {
		return gateway.ErrUserNotFound
	}

	pqerr, ok := (err).(*pq.Error)
	if !ok {
		return err
	}

	switch pqerr.Code {
	case InvalidTextRepresentation:
		return gateway.ErrUserNotFound
	case UniqueViolation:
		return gateway.ErrUserAlreadyExists
	}

	return err
}

func generateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

func formatAddress(addr string) string {
	rep := strings.NewReplacer(
		"pasaje", "Pje.",
		"Pasaje", "Pje.",
		"avenida", "Av.",
		"Avenida", "Av.",
		"esquina", "esq.",
		"Esquina", "esq.",
		"departamento", "depto.",
		"Departamento", "depto.",
		"º", "",
		"°", "",
	)

	addr = rep.Replace(nfcString(addr))
	if len(addr) > addrMaxLen {
		return addr[:addrMaxLen]
	}
	return addr
}
