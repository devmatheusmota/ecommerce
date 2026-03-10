package main

import (
	"log"
	"net/http"
	"os"

	"github.com/ecommerce/services/catalog/internal/database"
	"github.com/ecommerce/services/catalog/internal/router"
)

func main() {
	databaseConnection, err := database.Open()
	if err != nil {
		log.Fatal("database:", err)
	}
	defer databaseConnection.Close()

	migrateDatabase, err := database.Open()
	if err != nil {
		log.Fatal("database:", err)
	}
	if err := database.Migrate(migrateDatabase); err != nil {
		log.Fatal("migrate:", err)
	}
	migrateDatabase.Close()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	log.Printf("catalog service listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, router.New(databaseConnection)))
}
