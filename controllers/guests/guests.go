package guests

import (
	"database/sql"

	"github.com/jamesdavidyu/gender_reveal_service/cmd/model/types"
	"github.com/jamesdavidyu/gender_reveal_service/utils"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetGuestsByInviteeId(inviteeId string) (*types.Guests, error) {
	rows, err := s.db.Query("SELECT * FROM guests WHERE invitee_id = $1", inviteeId)
	if err != nil {
		return nil, err
	}

	guests := new(types.Guests)
	for rows.Next() {
		guests, err = utils.ScanRowIntoGuests(rows)
		if err != nil {
			return nil, err
		}
	}

	return guests, nil
}

func (s *Store) CreateGuests(guests types.Guests) error {
	_, err := s.db.Exec(
		`INSERT INTO guests (invitee_id, guests)
		VALUES ($1, $2)`,
		guests.InviteeId,
		guests.Guests,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) UpdateGuestsWithId(guests types.Guests) error {
	_, err := s.db.Exec(
		`UPDATE guests
		SET guests = $1
		WHERE invitee_id = $2`,
		guests.Guests,
		guests.InviteeId,
	)
	if err != nil {
		return err
	}

	return nil
}
