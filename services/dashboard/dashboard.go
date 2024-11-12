package dashboard

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jamesdavidyu/gender_reveal_service/cmd/model/types"
	"github.com/jamesdavidyu/gender_reveal_service/services/auth"
	"github.com/jamesdavidyu/gender_reveal_service/utils"
)

type Handler struct {
	store        types.DashboardStore
	inviteeStore types.InviteeStore
}

func NewHandler(store types.DashboardStore, inviteeStore types.InviteeStore) *Handler {
	return &Handler{store: store, inviteeStore: inviteeStore}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/dashboard/auth", auth.WithJWTAuth(h.handleGetRsvps, h.inviteeStore)).Methods("GET")
}

func (h *Handler) handleGetRsvps(w http.ResponseWriter, r *http.Request) {
	inviteeId := auth.GetInviteeIdFromContext(r.Context())

	role, err := h.inviteeStore.GetRoleById(inviteeId)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	if role.Role == "admin" {
		dashboard, err := h.store.GetRsvps()
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}

		utils.WriteJSON(w, http.StatusOK, dashboard)
	} else {
		auth.PermissionDenied(w)
		return
	}
}
