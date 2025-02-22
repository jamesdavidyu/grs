package newinvitee

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

func (s *Store) CreateNewInvitee(newInvitee types.NewInvitee) error {
	_, err := s.db.Exec(
		`INSERT INTO new_invitees (name, rsvp)
		VALUES ($1, $2)`,
		newInvitee.Name,
		newInvitee.Rsvp,
	)
	if err != nil {
		return err
	}

	return nil
}
