package routes

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	dashboardControllers "github.com/jamesdavidyu/gender_reveal_service/controllers/dashboard"
	guestsControllers "github.com/jamesdavidyu/gender_reveal_service/controllers/guests"
	inviteeControllers "github.com/jamesdavidyu/gender_reveal_service/controllers/invitees"
	newGuestsControllers "github.com/jamesdavidyu/gender_reveal_service/controllers/newGuests"
	newInviteeControllers "github.com/jamesdavidyu/gender_reveal_service/controllers/newInvitee"
	rsvpControllers "github.com/jamesdavidyu/gender_reveal_service/controllers/rsvp"
	dashboardServices "github.com/jamesdavidyu/gender_reveal_service/services/dashboard"
	guestsServices "github.com/jamesdavidyu/gender_reveal_service/services/guests"
	inviteeServices "github.com/jamesdavidyu/gender_reveal_service/services/invitees"
	newGuestsServices "github.com/jamesdavidyu/gender_reveal_service/services/newGuests"
	newInviteeServices "github.com/jamesdavidyu/gender_reveal_service/services/newInvitee"
	rsvpServices "github.com/jamesdavidyu/gender_reveal_service/services/rsvp"
	"github.com/jamesdavidyu/gender_reveal_service/utils"
	"github.com/joho/godotenv"
)

type APIServer struct {
	addr string
	db   *sql.DB
}

func NewAPIServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}
}

func (s *APIServer) Run() error {
	godotenv.Load()
	var Port = os.Getenv("PORT")

	router := mux.NewRouter()
	router.HandleFunc("/api/v1/status", getStatus()).Methods("GET")
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	inviteeStore := inviteeControllers.NewStore(s.db)
	inviteeHandler := inviteeServices.NewHandler(inviteeStore)
	inviteeHandler.RegisterRoutes(subrouter)

	rsvpStore := rsvpControllers.NewStore(s.db)
	rsvpHandler := rsvpServices.NewHandler(rsvpStore, inviteeStore)
	rsvpHandler.RegisterRoutes(subrouter)

	guestsStore := guestsControllers.NewStore(s.db)
	guestsHandler := guestsServices.NewHandler(guestsStore, inviteeStore)
	guestsHandler.RegisterRoutes(subrouter)

	dashboardStore := dashboardControllers.NewStore(s.db)
	dashboardHandler := dashboardServices.NewHandler(dashboardStore, inviteeStore)
	dashboardHandler.RegisterRoutes(subrouter)

	newInviteeStore := newInviteeControllers.NewStore(s.db)
	newInviteeHandler := newInviteeServices.NewHandler(newInviteeStore)
	newInviteeHandler.RegisterRoutes(subrouter)

	newGuestsStore := newGuestsControllers.NewStore(s.db)
	newGuestsHandler := newGuestsServices.NewHandler(newGuestsStore)
	newGuestsHandler.RegisterRoutes(subrouter)

	enhancedRouter := utils.EnableCORS(utils.JSONContentTypeMiddleware(router))

	log.Println("listening on", Port)
	return http.ListenAndServe(":"+Port, enhancedRouter)
}

func getStatus() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		okStatus := map[string]string{"status": "ok"}
		if err := json.NewEncoder(w).Encode(okStatus); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
