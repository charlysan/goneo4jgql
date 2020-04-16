package models

// Movie represents a movie
type Movie struct {
	UUID     string `json:"uuid" db:"uuid"`
	Title    string `json:"title" db:"title"`
	Tagline  string `json:"tagline,omitempty" db:"tagline"`
	Released int64  `json:"released" db:"released"`
}

// IsNode needed for gqlgen
func (i *Movie) IsNode() {}

// Person represents a person within a movie
type Person struct {
	UUID string  `json:"uuid" db:"uuid"`
	Name string  `json:"name" db:"name"`
	Role *string `json:"role" db:"role"`
	Born int64   `json:"born" db:"born"`
}

// IsNode needed for gqlgen
func (i *Person) IsNode() {}
