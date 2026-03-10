// Seed product images from the internet (Picsum Photos). Run from catalog dir with same env as the service:
//
//	cd services/catalog && go run ./cmd/seed-product-images
//
// Or with explicit DB: PG_HOST=localhost PG_DATABASE=catalog go run ./cmd/seed-product-images
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/ecommerce/services/catalog/internal/database"
	_ "github.com/lib/pq"
)

const baseURL = "https://picsum.photos/seed/%s/400/400"

// imageCount returns a deterministic count between 2 and 7 based on product id.
func imageCount(id string) int {
	var sum int
	for _, b := range id {
		sum += int(b)
	}
	return 2 + (sum % 6)
}

func main() {
	db, err := database.Open()
	if err != nil {
		log.Fatal("database:", err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT id FROM products")
	if err != nil {
		log.Fatal("query products:", err)
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			log.Fatal("scan:", err)
		}
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		log.Fatal("rows:", err)
	}

	if len(ids) == 0 {
		log.Println("no products found")
		return
	}

	log.Printf("updating %d products with 2–7 images each (Picsum Photos)", len(ids))

	for _, id := range ids {
		idNoDash := strings.ReplaceAll(id, "-", "")
		count := imageCount(id)
		images := make([]string, count)
		for i := range count {
			seed := fmt.Sprintf("%s%d", idNoDash, i)
			images[i] = fmt.Sprintf(baseURL, seed)
		}
		imagesJSON, err := json.Marshal(images)
		if err != nil {
			log.Printf("skip %s: marshal: %v", id, err)
			continue
		}
		_, err = db.Exec("UPDATE products SET images = $1 WHERE id = $2", imagesJSON, id)
		if err != nil {
			log.Printf("skip %s: update: %v", id, err)
			continue
		}
		log.Printf("updated %s -> %d images", id, count)
	}

	log.Printf("done: %d products updated", len(ids))
}
