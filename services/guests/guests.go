package guests

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/jamesdavidyu/gender_reveal_service/cmd/model/types"
	"github.com/jamesdavidyu/gender_reveal_service/services/auth"
	"github.com/jamesdavidyu/gender_reveal_service/utils"
)

type Handler struct {
	store        types.GuestsStore
	inviteeStore types.InviteeStore
}

func NewHandler(store types.GuestsStore, inviteeStore types.InviteeStore) *Handler {
	return &Handler{store: store, inviteeStore: inviteeStore}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/guests/auth", auth.WithJWTAuth(h.handleGetGuests, h.inviteeStore)).Methods("GET")
	router.HandleFunc("/guests/create/auth", auth.WithJWTAuth(h.handleCreateGuests, h.inviteeStore)).Methods("POST")
	router.HandleFunc("/guests/update/auth", auth.WithJWTAuth(h.handlePutGuests, h.inviteeStore)).Methods("PUT")
}

func (h *Handler) handleGetGuests(w http.ResponseWriter, r *http.Request) {
	inviteeId := auth.GetInviteeIdFromContext(r.Context())

	guests, err := h.store.GetGuestsByInviteeId(inviteeId)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("database error"))
		return
	}

	utils.WriteJSON(w, http.StatusOK, guests)
}

func (h *Handler) handleCreateGuests(w http.ResponseWriter, r *http.Request) {
	inviteeId := auth.GetInviteeIdFromContext(r.Context())
	var guests types.GuestsPayload

	if err := json.NewDecoder(r.Body).Decode(&guests); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("bad data"))
		return
	}

	if err := utils.Validate.Struct(guests); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid submission for %v", errors))
		return
	}

	err := h.store.CreateGuests(types.Guests{
		InviteeId: inviteeId,
		Guests:    guests.Guests,
	})
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("database error"))
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(guests)
}

func (h *Handler) handlePutGuests(w http.ResponseWriter, r *http.Request) {
	inviteeId := auth.GetInviteeIdFromContext(r.Context())
	var newGuests types.GuestsPayload

	if err := json.NewDecoder(r.Body).Decode(&newGuests); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("bad data"))
		return
	}

	err := h.store.UpdateGuestsWithId(types.Guests{
		InviteeId: inviteeId,
		Guests:    newGuests.Guests,
	})
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("database error"))
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newGuests)
}
