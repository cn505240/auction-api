package service

import (
	"auction-api/database"
	"fmt"

	"github.com/google/uuid"
)

type CreateBidRequest struct {
	Amount    float64 `json:"amount"`
	BidderID  string  `json:"bidder_id"`
	AuctionID string  `json:"auction_id"`
}

// TODO: Bidding with an associated max bid and ensuing automated bid handling

func (s *AuctionService) CreateBid(bidRequest *CreateBidRequest) (*database.Bid, error) {
	// Validate auction ID is not empty
	if bidRequest.AuctionID == "" {
		return nil, fmt.Errorf("auction ID is required")
	}

	// Validate auction ID is a valid UUID
	auctionID, err := uuid.Parse(bidRequest.AuctionID)
	if err != nil {
		return nil, fmt.Errorf("invalid auction ID format: %w", err)
	}

	// Check if the auction exists
	auction, err := s.repo.GetAuctionById(auctionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get auction: %w", err)
	}

	bidderID, err := uuid.Parse(bidRequest.BidderID)
	if err != nil {
		return nil, fmt.Errorf("invalid bidder ID: %w", err)
	}

	// Ensure bidder is not the seller
	if auction.SellerID == bidderID {
		return nil, fmt.Errorf("seller cannot bid on their own auction")
	}

	// Validate bid amount is positive
	if bidRequest.Amount <= 0 {
		return nil, fmt.Errorf("bid amount must be greater than zero")
	}

	// Get the current highest bid or use starting price
	var currentPrice float64 = auction.StartingPrice

	// Get the highest bid for the auction
	highestBid, err := s.repo.GetHighestBidForAuction(auctionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get highest bid for auction: %w", err)
	}

	if highestBid != nil {
		currentPrice = highestBid.Amount
	}

	// Calculate minimum bid increment based on current price
	var minIncrement float64
	switch {
	case currentPrice < 100:
		minIncrement = 1.0 // $1 increment for items under $100
	case currentPrice < 1000:
		minIncrement = 5.0 // $5 increment for items $100-$999
	case currentPrice < 5000:
		minIncrement = 10.0 // $10 increment for items $1000-$4999
	default:
		minIncrement = 50.0 // $50 increment for items $5000+
	}

	// Ensure bid meets minimum increment requirement
	if bidRequest.Amount < currentPrice+minIncrement {
		return nil, fmt.Errorf("bid amount must be at least %.2f (current price %.2f + minimum increment %.2f)",
			currentPrice+minIncrement, currentPrice, minIncrement)
	}

	bid := &database.Bid{
		AuctionID: auctionID,
		BidderID:  bidderID,
		Amount:    bidRequest.Amount,
	}

	// Save the bid to the database
	createdBid, err := s.repo.CreateBid(bid)
	if err != nil {
		return nil, fmt.Errorf("failed to create bid: %w", err)
	}

	return createdBid, nil

}
