package invitees

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

func (s *Store) GetInviteeById(id string) (*types.Invitees, error) {
	rows, err := s.db.Query("SELECT * FROM invitees WHERE id = $1", id)
	if err != nil {
		return nil, err
	}

	invitee := new(types.Invitees)
	for rows.Next() {
		invitee, err = utils.ScanRowIntoInvitee(rows)
		if err != nil {
			return nil, err
		}
	}

	return invitee, nil
}

func (s *Store) GetInviteeWithName(name string) (*types.Invitees, error) {
	rows, err := s.db.Query(
		`SELECT * FROM invitees
		WHERE name = $1`, name,
	)
	if err != nil {
		return nil, err
	}

	invitee := new(types.Invitees)
	for rows.Next() {
		invitee, err = utils.ScanRowIntoInvitee(rows)
		if err != nil {
			return nil, err
		}
	}

	return invitee, nil
}

func (s *Store) CreateInvitee(invitee types.Invitees) error {
	_, err := s.db.Exec(
		`INSERT INTO invitees (name, password)
		VALUES ($1, $2)`,
		invitee.Name,
		invitee.Password,
	)
	if err != nil {
		return err
	}

	return nil
}
