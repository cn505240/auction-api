package service

import (
	"auction-api/database"
	"fmt"

	"github.com/google/uuid"
)

// CreateAuctionRequest represents the data structure for auction creation from API
type CreateAuctionRequest struct {
	Title         string  `json:"title"`
	Description   string  `json:"description"`
	StartingPrice float64 `json:"starting_price"`
	SellerID      string  `json:"seller_id"`
}

// CreateAuction creates a new auction from API request data
func (s *AuctionService) CreateAuction(req *CreateAuctionRequest) (*database.Auction, error) {
	if req.Title == "" {
		return nil, fmt.Errorf("auction title cannot be empty")
	}

	if req.StartingPrice <= 0 {
		return nil, fmt.Errorf("starting price must be greater than zero")
	}

	if req.SellerID == "" {
		return nil, fmt.Errorf("seller ID cannot be empty")
	}

	// Validate that the seller ID is a valid UUID
	sellerID, err := uuid.Parse(req.SellerID)
	if err != nil {
		return nil, fmt.Errorf("invalid seller ID format: %w", err)
	}

	// Check if the seller exists in the database
	_, err = s.repo.GetUserById(sellerID)
	if err != nil {
		return nil, fmt.Errorf("invalid seller ID: %w", err)
	}

	// Create auction entity from request
	auction := &database.Auction{
		Title:         req.Title,
		Description:   req.Description,
		StartingPrice: req.StartingPrice,
		SellerID:      sellerID,
	}

	createdAuction, err := s.repo.CreateAuction(auction)
	if err != nil {
		return nil, fmt.Errorf("failed to create auction: %w", err)
	}

	return createdAuction, nil
}

// TODO: Add GetAuction and ListAuctions methods
