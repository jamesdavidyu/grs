package routes

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	inviteeControllers "github.com/jamesdavidyu/gender_reveal_service/controllers/invitees"
	inviteeServices "github.com/jamesdavidyu/gender_reveal_service/services/invitees"
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
