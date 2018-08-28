package users

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/lib/pq"
	"github.com/rodrwan/gateway"
	"github.com/rodrwan/gateway/postgres"
)

const (
	addrMaxLen = 30
)

// UserStore implements gateway.UserStore interface with postgres as backend
type UserStore struct {
	store postgres.SQLExecutor
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
	if opt.Email != "" {
		q = q.Where("email = ?", opt.Email)
	}

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

// UsersWithAddress ...
func (us *UserStore) UsersWithAddress() ([]*gateway.User, error) {
	rows, err := us.store.Queryx(
		`select
	users.*,
	addresses.user_id "addresses.user_id",
	addresses.address_line "addresses.address_line",
	addresses.city "addresses.city",
	addresses.locality "addresses.locality",
	addresses.administrative_area_level_1 "addresses.administrative_area_level_1",
	addresses.country "addresses.country",
	addresses.postal_code "addresses.postal_code"
	from users inner join addresses on addresses.user_id = users.id
`,
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

func userError(err error) error {
	if err == sql.ErrNoRows {
		return gateway.ErrUserNotFound
	}

	pqerr, ok := (err).(*pq.Error)
	if !ok {
		return err
	}

	switch pqerr.Code {
	case postgres.InvalidTextRepresentation:
		return gateway.ErrUserNotFound
	case postgres.UniqueViolation:
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

	addr = rep.Replace(postgres.NFCString(addr))
	if len(addr) > addrMaxLen {
		return addr[:addrMaxLen]
	}
	return addr
}
