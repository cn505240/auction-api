package database

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Bid represents a bid in the system
type Bid struct {
	ID        uuid.UUID `json:"id"`
	AuctionID uuid.UUID `json:"auction_id"`
	BidderID  uuid.UUID `json:"bidder_id"`
	Amount    float64   `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateBid inserts a new bid into the database and returns the created bid
func (db *DB) CreateBid(bid *Bid) (*Bid, error) {
	now := time.Now()

	bid.ID = uuid.New()
	bid.CreatedAt = now
	bid.UpdatedAt = now

	query := `
		INSERT INTO bids (id, auction_id, bidder_id, amount, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := db.conn.Exec(query,
		bid.ID,
		bid.AuctionID,
		bid.BidderID,
		bid.Amount,
		bid.CreatedAt,
		bid.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create bid: %w", err)
	}

	return bid, nil
}

// GetHighestBidForAuction retrieves the highest bid for a specific auction
func (db *DB) GetHighestBidForAuction(auctionID uuid.UUID) (*Bid, error) {
	query := `
		SELECT id, auction_id, bidder_id, amount, created_at, updated_at
		FROM bids
		WHERE auction_id = $1
		ORDER BY created_at DESC
		LIMIT 1
	`

	var bid Bid
	err := db.conn.QueryRow(query, auctionID).Scan(
		&bid.ID,
		&bid.AuctionID,
		&bid.BidderID,
		&bid.Amount,
		&bid.CreatedAt,
		&bid.UpdatedAt,
	)

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, nil // No bids found for this auction
		}
		return nil, fmt.Errorf("failed to get highest bid for auction: %w", err)
	}

	return &bid, nil
}
