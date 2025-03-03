package newinvitee

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
	store types.NewInviteeStore
}

func NewHandler(store types.NewInviteeStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/newinvitee", h.handleCreateNewInvitee).Methods("POST")
}

func (h *Handler) handleCreateNewInvitee(w http.ResponseWriter, r *http.Request) {
	var newInvitee types.NewInviteePayload

	if err := json.NewDecoder(r.Body).Decode(&newInvitee); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("bad data"))
	}

	if err := utils.Validate.Struct(newInvitee); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid submission for %v", errors))
	}

	err := h.store.CreateNewInvitee(types.NewInvitee{
		Name: newInvitee.Name,
	})
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("database error"))
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newInvitee)
}
