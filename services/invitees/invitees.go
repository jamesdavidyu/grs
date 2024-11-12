package invitees

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/jamesdavidyu/gender_reveal_service/cmd/model/types"
	"github.com/jamesdavidyu/gender_reveal_service/config"
	"github.com/jamesdavidyu/gender_reveal_service/services/auth"
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
	router.HandleFunc("/auth/register", auth.WithJWTAuth(h.handleRegister, h.store)).Methods("POST")
	router.HandleFunc("/auth/login", h.handleLogin).Methods("POST")
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
		return
	}

	inviteeId := auth.GetInviteeIdFromContext(r.Context())

	role, err := h.store.GetRoleById(inviteeId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("bad data"))
		return
	}

	if role.Role == "admin" {
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
	} else {
		auth.PermissionDenied(w)
		return
	}
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	var login types.Register
	if err := json.NewDecoder(r.Body).Decode(&login); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("bad data"))
		return
	}

	if err := utils.Validate.Struct(login); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid submission for %v", errors))
		return
	}

	invitee, err := h.store.GetInviteeWithName(login.Name)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("database error"))
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(invitee.Password), []byte(login.Password)) != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("database error"))
		return
	} else {
		secret := []byte(config.Envs.JWTSecret)
		token, err := auth.CreateJWT(secret, invitee.Id)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("database error"))
			return
		}

		json.NewEncoder(w).Encode(map[string]any{
			"token":     token,
			"inviteeId": invitee.Id,
			"name":      invitee.Name,
		})
	}
}
