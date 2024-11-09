package rsvp

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
	store        types.RsvpStore
	inviteeStore types.InviteeStore
}

func NewHandler(store types.RsvpStore, inviteeStore types.InviteeStore) *Handler {
	return &Handler{store: store, inviteeStore: inviteeStore}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/rsvp/auth", auth.WithJWTAuth(h.handleGetRsvp, h.inviteeStore)).Methods("GET")
	router.HandleFunc("/rsvp/create/auth", auth.WithJWTAuth(h.handleCreateRsvp, h.inviteeStore)).Methods("POST")
	router.HandleFunc("/rsvp/update/auth", auth.WithJWTAuth(h.handlePutRsvp, h.inviteeStore)).Methods("PUT")
}

func (h *Handler) handleGetRsvp(w http.ResponseWriter, r *http.Request) {
	inviteeId := auth.GetInviteeIdFromContext(r.Context())

	rsvp, err := h.store.GetRsvpByInviteeId(inviteeId)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("database error"))
		return
	}

	utils.WriteJSON(w, http.StatusOK, rsvp)
}

func (h *Handler) handleCreateRsvp(w http.ResponseWriter, r *http.Request) {
	inviteeId := auth.GetInviteeIdFromContext(r.Context())
	var rsvp types.CreateRsvpPayload

	if err := json.NewDecoder(r.Body).Decode(&rsvp); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("bad data"))
		return
	}

	if err := utils.Validate.Struct(rsvp); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid submission for %v", errors))
		return
	}

	err := h.store.CreateRsvp(types.Rsvp{
		InviteeId: inviteeId,
		Rsvp:      rsvp.Rsvp,
	})
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(rsvp)
}

func (h *Handler) handlePutRsvp(w http.ResponseWriter, r *http.Request) {
	inviteeId := auth.GetInviteeIdFromContext(r.Context())
	var newRsvp types.CreateRsvpPayload

	if err := json.NewDecoder(r.Body).Decode(&newRsvp); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("bad data"))
		return
	}

	err := h.store.UpdateRsvpWithId(types.Rsvp{
		InviteeId: inviteeId,
		Rsvp:      newRsvp.Rsvp,
	})
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("database error"))
		return
	}
}
