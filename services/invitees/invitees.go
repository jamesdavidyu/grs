package invitees

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/jamesdavidyu/gender_reveal_service/cmd/model/types"
	"github.com/jamesdavidyu/gender_reveal_service/utils"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	store types.InviteeStore
}

func NewHandler(store types.InviteeStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/auth/register", h.handleRegister).Methods("POST")
}

func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	var register types.Register
	if err := json.NewDecoder(r.Body).Decode(&register); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("bad data"))
		return
	}

	if err := utils.Validate.Struct(register); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid submission for %v", errors))
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(register.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("database error"))
		return
	}

	err = h.store.CreateInvitee(types.Invitees{
		Name:     register.Name,
		Password: string(hashedPassword),
	})
	if err != nil {
		utils.WriteError(w, http.StatusAlreadyReported, fmt.Errorf("already entered"))
		return
	}
}
