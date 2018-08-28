package addresses

import (
	"database/sql"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/lib/pq"
	"github.com/rodrwan/gateway"
	"github.com/rodrwan/gateway/postgres"
)

// Store implements gateway.AddressStore interface with postgres as backend
type Store struct {
	Store postgres.SQLExecutor
}

// Find search an address by user_id .
func (as *Store) Find(userID string) (*gateway.Address, error) {
	row := as.Store.QueryRowx(
		"select * from addresses where user_id = $1",
		userID,
	)
	var a gateway.Address
	if err := row.StructScan(&a); err != nil {
		return nil, addressError(err)
	}
	return &a, nil
}

// Create creates a new address
func (as *Store) Create(a *gateway.Address) error {
	sql, args, err := squirrel.Insert("addresses").
		Columns(
			"user_id", "city", "address_line", "locality", "administrative_area_level_1",
			"country", "postal_code",
		).
		Values(
			a.UserID, a.City, a.AddressLine, a.Locality, a.AdministrativeAreaLevel1,
			a.Country, a.PostalCode,
		).
		Suffix("returning *").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return addressError(err)
	}

	row := as.Store.QueryRowx(sql, args...)
	if err := row.StructScan(a); err != nil {
		return addressError(err)
	}

	return nil
}

// Update update the given address
func (as *Store) Update(a *gateway.Address) error {
	sql, args, err := squirrel.Update("addresses").
		Set("city", a.City).
		Set("address_line", a.AddressLine).
		Set("locality", a.Locality).
		Set("administrative_area_level_1", a.AdministrativeAreaLevel1).
		Set("country", a.Country).
		Set("postal_code", a.PostalCode).
		Where("user_id = ?", a.UserID).
		Suffix("returning *").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	row := as.Store.QueryRowx(sql, args...)
	if err := row.StructScan(a); err != nil {
		return addressError(err)
	}

	return nil
}

func addressError(err error) error {
	if err == sql.ErrNoRows {
		return err
	}

	pqerr, ok := (err).(*pq.Error)
	if !ok {
		return err
	}

	switch pqerr.Code {
	case postgres.InvalidTextRepresentation:
		if strings.Contains(pqerr.Message, "user_requirement") {
			return err
		}
		if strings.Contains(pqerr.Message, "card_provider") {
			return err
		}
		return err
	case postgres.UniqueViolation:
		return err
	}

	return err
}
