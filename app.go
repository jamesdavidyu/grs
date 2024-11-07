package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/jamesdavidyu/gender_reveal_service/cmd/model/db"
	"github.com/jamesdavidyu/gender_reveal_service/routes"
	"github.com/joho/godotenv"
)

// go:embed templates/*
// var resources embed.FS

// var t = template.Must(template.ParseFS(resources, "templates/*"))

func main() {
	godotenv.Load()
	var Port = os.Getenv("PORT")

	db, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}

	initStorage(db)

	server := routes.NewAPIServer(":"+Port, db)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}

	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	data := map[string]string{
	// 		"Region": os.Getenv("FLY_REGION"),
	// 	}

	// 	t.ExecuteTemplate(w, "index.html.tmpl", data)
	// })
}

func initStorage(db *sql.DB) {
	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("DB: Successfully connected!")
}
