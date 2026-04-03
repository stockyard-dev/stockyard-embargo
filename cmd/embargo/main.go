package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/stockyard-dev/stockyard-embargo/internal/server"
	"github.com/stockyard-dev/stockyard-embargo/internal/store"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "9150"
	}
	dataDir := os.Getenv("DATA_DIR")
	if dataDir == "" {
		dataDir = "./embargo-data"
	}

	db, err := store.Open(dataDir)
	if err != nil {
		log.Fatalf("embargo: open database: %v", err)
	}
	defer db.Close()

	srv := server.New(db, server.DefaultLimits())

	fmt.Printf("\n  Embargo — Self-hosted content embargo manager\n")
	fmt.Printf("  ─────────────────────────────────\n")
	fmt.Printf("  Dashboard:  http://localhost:%s/ui\n", port)
	fmt.Printf("  API:        http://localhost:%s/api\n", port)
	fmt.Printf("  Data:       %s\n", dataDir)
	fmt.Printf("  ─────────────────────────────────\n\n")

	log.Printf("embargo: listening on :%s", port)
	if err := http.ListenAndServe(":"+port, srv); err != nil {
		log.Fatalf("embargo: %v", err)
	}
}
