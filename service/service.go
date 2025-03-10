package service

import (
	"auction-api/database"

	"github.com/google/uuid"
)

// AuctionRepository defines the interface for auction database operations
type AuctionRepository interface {
	GetUserById(id uuid.UUID) (*database.User, error)
	CreateAuction(auction *database.Auction) (*database.Auction, error)
	GetAuctionById(id uuid.UUID) (*database.Auction, error)
	CreateBid(bid *database.Bid) (*database.Bid, error)
	GetHighestBidForAuction(auctionID uuid.UUID) (*database.Bid, error)
}

type AuctionService struct {
	repo AuctionRepository
}

// NewAuctionService creates a new auction service
func NewAuctionService(repo AuctionRepository) *AuctionService {
	return &AuctionService{
		repo: repo,
	}
}
