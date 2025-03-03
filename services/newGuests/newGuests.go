package newguests

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/jamesdavidyu/gender_reveal_service/cmd/model/types"
	"github.com/jamesdavidyu/gender_reveal_service/utils"
)

type Handler struct {
	store types.NewGuestsStore
}

func NewHandler(store types.NewGuestsStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/newguests", h.handleCreateNewGuests).Methods("POST")
}

func (h *Handler) handleCreateNewGuests(w http.ResponseWriter, r *http.Request) {
	var newGuests types.NewGuestsPayload

	if err := json.NewDecoder(r.Body).Decode(&newGuests); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("bad data"))
	}

	if err := utils.Validate.Struct(newGuests); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid submission for %v", errors))
	}

	err := h.store.CreateNewGuests(types.NewGuests{
		Name:   newGuests.Name,
		Rsvp:   newGuests.Rsvp,
		Guests: newGuests.Guests,
	})
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("database error"))
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newGuests)
}
