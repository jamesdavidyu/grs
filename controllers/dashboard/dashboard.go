package dashboard

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

func (s *Store) GetRsvps() ([]types.Dashboard, error) {
	rows, err := s.db.Query(
		`SELECT 
			r.id, 
			i.name, 
			r.rsvp, 
			(CASE 
				WHEN g.guests IS NULL THEN 'none'
				ELSE g.guests
			END) AS guests 
		FROM rsvp r
		LEFT JOIN invitees i ON i.id = r.invitee_id
		LEFT JOIN guests g ON g.invitee_id = r.invitee_id`,
	)
	if err != nil {
		return nil, err
	}

	dashboards := make([]types.Dashboard, 0)
	for rows.Next() {
		dashboard, err := utils.ScanRowIntoDashboard(rows)
		if err != nil {
			return nil, err
		}
		dashboards = append(dashboards, *dashboard)
	}

	return dashboards, nil
}
