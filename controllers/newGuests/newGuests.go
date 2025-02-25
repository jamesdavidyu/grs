package newguests

import (
	"database/sql"

	"github.com/jamesdavidyu/gender_reveal_service/cmd/model/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CreateNewGuests(newGuests types.NewGuests) error {
	_, err := s.db.Exec(
		`INSERT INTO new_guests (name, guests)
		VALUES ($1, $2)`,
		newGuests.Name,
		newGuests.Guests,
	)
	if err != nil {
		return err
	}

	return nil
}
