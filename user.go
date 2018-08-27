package gateway

import (
	"errors"
	"time"
)

// User represents a user this hold user information returned by UserService/
type User struct {
	ID    string `json:"id,omitempty" db:"id"`
	Email string `json:"email,omitempty" db:"email"`

	FirstName string    `json:"first_name,omitempty" db:"first_name"`
	LastName  string    `json:"last_name,omitempty" db:"last_name"`
	Phone     string    `json:"phone,omitempty" db:"phone"`
	Birthdate time.Time `json:"birthdate,omitempty" db:"birthdate"`

	Address *Address `json:"address,omitempty" db:"-"`
}

// UserService error definitions.
var (
	ErrUserNotFound       = errors.New("user: user not found")
	ErrUserAlreadyExists  = errors.New("user: user already exists")
	ErrUserInvalidAddress = errors.New("user: invalid address fields")
	ErrAddressNotFound    = errors.New("user: address not found")
)

// UserQueryOptions represents query options to get users.
type UserQueryOptions struct {
	ID        string `q:"id"`
	Email     string `q:"email"`
	FirstName string `q:"first_name"`
	LastName  string `q:"last_name"`
}

// UserQueryOption provides an option to set User query options.
type UserQueryOption func(*UserQueryOptions)

// SetUserQueryOptions method to set parameters into UserQueryOptions struct.
func SetUserQueryOptions(src *UserQueryOptions) UserQueryOption {
	return func(dst *UserQueryOptions) {
		*dst = *src
	}
}

// UserService define service to be implemented by postgres service.
type UserService interface {
	Get(...UserQueryOption) (*User, error)
	Select(...UserQueryOption) ([]*User, error)
	All() ([]*User, error)
}
