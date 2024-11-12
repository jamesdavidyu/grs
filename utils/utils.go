package utils

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/jamesdavidyu/gender_reveal_service/cmd/model/types"
)

var Validate = validator.New()

func EnableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*") // TODO: change this to genderrevealparty.vercel.app?
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func JSONContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(v)
}

func WriteError(w http.ResponseWriter, status int, err error) {
	WriteJSON(w, status, map[string]string{"error": err.Error()})
}

func ScanRowIntoInvitee(rows *sql.Rows) (*types.Invitees, error) {
	invitee := new(types.Invitees)

	err := rows.Scan(
		&invitee.Id,
		&invitee.Name,
		&invitee.Password,
		&invitee.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return invitee, nil
}

func ScanRowIntoRsvp(rows *sql.Rows) (*types.Rsvp, error) {
	rsvp := new(types.Rsvp)

	err := rows.Scan(
		&rsvp.Id,
		&rsvp.InviteeId,
		&rsvp.Rsvp,
		&rsvp.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return rsvp, nil
}

func ScanRowIntoGuests(rows *sql.Rows) (*types.Guests, error) {
	guests := new(types.Guests)

	err := rows.Scan(
		&guests.Id,
		&guests.InviteeId,
		&guests.Guests,
		&guests.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return guests, nil
}

func ScanRowIntoDashboard(rows *sql.Rows) (*types.Dashboard, error) {
	dashboard := new(types.Dashboard)

	err := rows.Scan(
		&dashboard.Id,
		&dashboard.Name,
		&dashboard.Rsvp,
		&dashboard.Guests,
	)
	if err != nil {
		return nil, err
	}

	return dashboard, nil
}
