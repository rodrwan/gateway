package postgres

import (
	"unicode"

	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

// Postgres error codes
const (
	UniqueViolation           pq.ErrorCode = "23505"
	ForeignKeyViolation       pq.ErrorCode = "23503"
	InvalidTextRepresentation pq.ErrorCode = "22P02"
	NoDataFound               pq.ErrorCode = "P0002"
)

// Datastore store data in db use as a backend using postgres.
type Datastore struct {
	*sqlx.DB
}

// NewDatastore create a new Datastore instance
func NewDatastore(dsn string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(4)
	db.SetMaxOpenConns(16)
	return db, nil
}

func nfcString(nfd string) (nfc string) {
	isMn := func(r rune) bool {
		return unicode.Is(unicode.Mn, r) // Mn: nonspacing marks
	}

	t := transform.Chain(norm.NFD, transform.RemoveFunc(isMn), norm.NFC)
	nfc, _, _ = transform.String(t, nfd)
	return
}

// SQLExecutor provides an abstraction layer over a SQL executor.
type SQLExecutor interface {
	sqlx.Queryer
	sqlx.Execer
}
