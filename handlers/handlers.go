package handlers

import (
	"auction-api/database"
	"auction-api/service"
)

// AuctionService defines the interface for auction service operations
type AuctionService interface {
	CreateAuction(req *service.CreateAuctionRequest) (*database.Auction, error)
	CreateBid(req *service.CreateBidRequest) (*database.Bid, error)
}

// Handlers struct holds dependencies for the HTTP handlers
type Handlers struct {
	service AuctionService
}

// NewHandlers creates a new Handlers instance
func NewHandlers(service AuctionService) *Handlers {
	return &Handlers{
		service: service,
	}
}
