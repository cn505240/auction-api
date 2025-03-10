package main

import (
	"auction-api/database"
	"auction-api/handlers"
	"auction-api/service"
	"database/sql"
	"log"
	"net/http"
	"regexp"
)

func main() {

	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost/postgres?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := database.RunMigrations(db); err != nil {
		log.Fatal(err)
	}

	// Initialize repository, service, and handlers
	repo, err := database.New("postgres://postgres:postgres@localhost/postgres?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	auctionService := service.NewAuctionService(repo)
	h := handlers.NewHandlers(auctionService)

	// Define routes
	http.HandleFunc("/auctions", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			h.CreateAuction(w, r)
			return
		}
		http.NotFound(w, r)
	})

	// Handle bids with regex to extract auction ID
	bidPattern := regexp.MustCompile(`^/auctions/([^/]+)/bids$`)
	http.HandleFunc("/auctions/", func(w http.ResponseWriter, r *http.Request) {
		matches := bidPattern.FindStringSubmatch(r.URL.Path)
		if matches != nil && r.Method == http.MethodPost {
			// URL path matches the pattern for creating a bid
			auctionID := matches[1] // Extract the auction ID from the regex match
			h.CreateBid(w, r, auctionID)
			return
		}
		http.NotFound(w, r)
	})

	// Start the server
	log.Println("Starting Auction API on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Server error:", err)
	}
}
