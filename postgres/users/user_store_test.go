package users

import (
	"encoding/json"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/rodrwan/gateway"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

var dummyUsers = []byte(`
[
    {
        "id": "1",
        "first_name": "louane",
        "last_name": "vidal",
        "email": "louane.vidal@example.com",
        "bithdate": "1965-12-20T13:32:15Z",
        "phone": "02-62-35-18-98"
    },
    {
        "id": "2",
        "first_name": "noel",
        "last_name": "peixoto",
        "email": "noel.peixoto@example.com",
        "birthdate": "1954-09-27T05:27:22Z",
        "phone": "(27) 2001-0083"
    },
    {
        "id": "3",
        "first_name": "manuel",
        "last_name": "lorenzo",
        "email": "manuel.lorenzo@example.com",
        "bithdate": "1949-05-26T09:55:03Z",
        "phone": "936-865-442"
    }
]`)

func generateRows() *sqlmock.Rows {
	rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "email", "phone", "birthdate"})
	users := []*gateway.User{}

	if err := json.Unmarshal(dummyUsers, &users); err != nil {
		panic(err)
	}

	for _, user := range users {
		rows.AddRow(user.ID, user.FirstName, user.LastName, user.Email, user.Phone, user.Birthdate)
	}

	return rows
}

func generateRow() *sqlmock.Rows {
	rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "email", "phone", "birthdate"})
	users := []*gateway.User{}

	if err := json.Unmarshal(dummyUsers, &users); err != nil {
		panic(err)
	}

	user := users[0]
	rows.AddRow(user.ID, user.FirstName, user.LastName, user.Email, user.Phone, user.Birthdate)

	return rows
}

func TestGet_UserStore(t *testing.T) {
	driverName := "sqlmock"
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	store := sqlx.NewDb(mockDB, driverName)
	us := UserStore{store}

	rows := generateRows()
	email := "louane.vidal@example.com"
	mock.ExpectQuery("SELECT (.+) FROM users WHERE email = ?").
		WithArgs(email).
		WillReturnRows(rows)

	opts := gateway.SetUserQueryOptions(&gateway.UserQueryOptions{
		Email: email,
	})
	u, err := us.Get(opts)
	if err != nil {
		t.Fatalf("an error '%s' was not expected", err)
	}

	firstName := "louane"
	if u.FirstName != firstName {
		t.Errorf("expected mocked first_name to be '%s', but got '%s' instead", firstName, u.FirstName)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestSelect_UserStore(t *testing.T) {
	driverName := "sqlmock"
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	store := sqlx.NewDb(mockDB, driverName)

	rows := generateRow()
	email := "louane.vidal@example.com"
	mock.ExpectQuery("SELECT (.+) FROM users WHERE deleted_at is null AND email = ?").
		WithArgs(email).
		WillReturnRows(rows)

	us := UserStore{store}

	opts := gateway.SetUserQueryOptions(&gateway.UserQueryOptions{
		Email: email,
	})
	uu, err := us.Select(opts)
	if err != nil {
		t.Fatalf("an error '%s' was not expected", err)
	}

	expectedLength := 1
	if len(uu) != expectedLength {
		t.Fatalf("expected mocked length to be '%d', but got '%d' instead", expectedLength, len(uu))
	}

	firstName := "louane"
	if uu[0].FirstName != firstName {
		t.Errorf("expected mocked first_name to be '%s', but got '%s' instead", firstName, uu[0].FirstName)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestAll_UserStore(t *testing.T) {
	driverName := "sqlmock"
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	store := sqlx.NewDb(mockDB, driverName)
	us := UserStore{store}

	rows := generateRows()

	mock.ExpectQuery("select (.+) from users").WillReturnRows(rows)

	uu, err := us.All()
	if err != nil {
		t.Fatalf("an error '%s' was not expected", err)
	}

	expectedLength := 3
	if len(uu) != expectedLength {
		t.Fatalf("expected mocked length to be '%d', but got '%d' instead", expectedLength, len(uu))
	}
	firstName := "louane"
	if uu[0].FirstName != firstName {
		t.Errorf("expected mocked first_name to be '%s', but got '%s' instead", firstName, uu[0].FirstName)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
