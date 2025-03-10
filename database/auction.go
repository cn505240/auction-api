package database

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Auction represents an auction in the system
type Auction struct {
	ID            uuid.UUID `json:"id"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	StartingPrice float64   `json:"starting_price"`
	SellerID      uuid.UUID `json:"seller_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// CreateAuction inserts a new auction into the database
func (db *DB) CreateAuction(auction *Auction) (*Auction, error) {
	query := `
		INSERT INTO auctions (id, title, description, starting_price, seller_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, title, description, starting_price, seller_id, created_at, updated_at
	`

	now := time.Now()
	id := uuid.New()

	var result Auction
	err := db.conn.QueryRow(query,
		id,
		auction.Title,
		auction.Description,
		auction.StartingPrice,
		auction.SellerID,
		now,
		now,
	).Scan(
		&result.ID,
		&result.Title,
		&result.Description,
		&result.StartingPrice,
		&result.SellerID,
		&result.CreatedAt,
		&result.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create auction: %w", err)
	}

	return &result, nil
}

// GetAuctionById retrieves an auction from the database by its ID
func (db *DB) GetAuctionById(id uuid.UUID) (*Auction, error) {
	query := `
		SELECT id, title, description, starting_price, seller_id, created_at, updated_at
		FROM auctions
		WHERE id = $1
	`

	var auction Auction
	err := db.conn.QueryRow(query, id).Scan(
		&auction.ID,
		&auction.Title,
		&auction.Description,
		&auction.StartingPrice,
		&auction.SellerID,
		&auction.CreatedAt,
		&auction.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get auction by id: %w", err)
	}

	return &auction, nil
}
