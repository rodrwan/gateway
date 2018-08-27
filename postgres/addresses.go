package postgres

import (
	"database/sql"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/lib/pq"
	"github.com/rodrwan/gateway"
)

// AddressStore implements gateway.AddressStore interface with postgres as backend
type AddressStore struct {
	store SQLExecutor
}

// Find search an address by user_id .
func (as *AddressStore) Find(userID string) (*gateway.Address, error) {
	row := as.store.QueryRowx(
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
func (as *AddressStore) Create(a *gateway.Address) error {
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

	row := as.store.QueryRowx(sql, args...)
	if err := row.StructScan(a); err != nil {
		return addressError(err)
	}

	return nil
}

// Update update the given address
func (as *AddressStore) Update(a *gateway.Address) error {
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

	row := as.store.QueryRowx(sql, args...)
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
	case InvalidTextRepresentation:
		if strings.Contains(pqerr.Message, "user_requirement") {
			return err
		}
		if strings.Contains(pqerr.Message, "card_provider") {
			return err
		}
		return err
	case UniqueViolation:
		return err
	}

	return err
}
