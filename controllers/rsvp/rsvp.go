package rsvp

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

func (s *Store) GetRsvpByInviteeId(inviteeId string) (*types.Rsvp, error) {
	rows, err := s.db.Query("SELECT * FROM rsvp WHERE invitee_id = $1", inviteeId)
	if err != nil {
		return nil, err
	}

	rsvp := new(types.Rsvp)
	for rows.Next() {
		rsvp, err = utils.ScanRowIntoRsvp(rows)
		if err != nil {
			return nil, err
		}
	}

	return rsvp, nil
}

func (s *Store) CreateRsvp(rsvp types.Rsvp) error {
	_, err := s.db.Exec(
		`INSERT INTO rsvp (invitee_id, rsvp)
		VALUES ($1, $2)`,
		rsvp.InviteeId,
		rsvp.Rsvp,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) UpdateRsvpWithId(rsvp types.Rsvp) error {
	_, err := s.db.Exec(
		`UPDATE rsvp
		SET rsvp = $1
		WHERE invitee_id = $2`,
		rsvp.Rsvp,
		rsvp.InviteeId,
	)
	if err != nil {
		return err
	}

	return nil
}
