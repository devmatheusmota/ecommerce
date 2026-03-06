package main

import (
	"log"
	"net/http"
	"os"

	"github.com/ecommerce/services/users/internal/database"
	"github.com/ecommerce/services/users/internal/router"
)

func main() {
	db, err := database.Open()
	if err != nil {
		log.Fatal("database:", err)
	}
	defer db.Close()

	migrateDB, err := database.Open()
	if err != nil {
		log.Fatal("database:", err)
	}
	if err := database.Migrate(migrateDB); err != nil {
		log.Fatal("migrate:", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("users service listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, router.New(db)))
}
