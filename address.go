package gateway

// Address represents a user address this hold adderss information returned by AddressService
type Address struct {
	UserID                   string `json:"user_id,omitempty" db:"user_id"`
	AddressLine              string `json:"address_line,omitempty" db:"address_line"`
	City                     string `json:"city,omitempty" db:"city"`
	Locality                 string `json:"locality,omitempty" db:"locality"`
	AdministrativeAreaLevel1 string `json:"administrative_area_level_1,omitempty" db:"administrative_area_level_1"`
	Country                  string `json:"country,omitempty" db:"country"`
	PostalCode               int    `json:"postal_code,omitempty" db:"postal_code"`
}

// AddressService define service to be implemented by postgres service.
type AddressService interface{}
