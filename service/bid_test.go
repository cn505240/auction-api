package service

import (
	"auction-api/database"
	"errors"
	"testing"
	"time"

	mockservice "auction-api/mocks/auction-api/service"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type BidServiceTestSuite struct {
	suite.Suite
	mockRepo *mockservice.AuctionRepository
	service  *AuctionService
}

func (suite *BidServiceTestSuite) SetupTest() {
	suite.mockRepo = &mockservice.AuctionRepository{
		Mock: mock.Mock{},
	}
	suite.service = NewAuctionService(suite.mockRepo)
}

func TestBidServiceSuite(t *testing.T) {
	suite.Run(t, new(BidServiceTestSuite))
}

func (suite *BidServiceTestSuite) TestCreateBid_Success() {
	// Arrange
	auctionID := uuid.New()
	bidderID := uuid.New()
	sellerID := uuid.New()

	req := &CreateBidRequest{
		Amount:    150.0,
		BidderID:  bidderID.String(),
		AuctionID: auctionID.String(),
	}

	mockAuction := &database.Auction{
		ID:            auctionID,
		Title:         "Test Auction",
		Description:   "This is a test auction",
		StartingPrice: 100.0,
		SellerID:      sellerID,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	expectedBid := &database.Bid{
		AuctionID: auctionID,
		BidderID:  bidderID,
		Amount:    req.Amount,
	}

	returnedBid := &database.Bid{
		ID:        uuid.New(),
		AuctionID: auctionID,
		BidderID:  bidderID,
		Amount:    req.Amount,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Mock expectations
	suite.mockRepo.On("GetAuctionById", auctionID).Return(mockAuction, nil)
	suite.mockRepo.On("GetHighestBidForAuction", auctionID).Return(nil, nil)
	suite.mockRepo.On("CreateBid", mock.MatchedBy(func(b *database.Bid) bool {
		return b.AuctionID == expectedBid.AuctionID &&
			b.BidderID == expectedBid.BidderID &&
			b.Amount == expectedBid.Amount
	})).Return(returnedBid, nil)

	// Act
	result, err := suite.service.CreateBid(req)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), returnedBid, result)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *BidServiceTestSuite) TestCreateBid_EmptyAuctionID() {
	// Arrange
	req := &CreateBidRequest{
		Amount:    150.0,
		BidderID:  uuid.New().String(),
		AuctionID: "",
	}

	// Act
	result, err := suite.service.CreateBid(req)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Contains(suite.T(), err.Error(), "auction ID is required")
	suite.mockRepo.AssertNotCalled(suite.T(), "GetAuctionById")
	suite.mockRepo.AssertNotCalled(suite.T(), "CreateBid")
}

func (suite *BidServiceTestSuite) TestCreateBid_InvalidAuctionID() {
	// Arrange
	req := &CreateBidRequest{
		Amount:    150.0,
		BidderID:  uuid.New().String(),
		AuctionID: "invalid-uuid",
	}

	// Act
	result, err := suite.service.CreateBid(req)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Contains(suite.T(), err.Error(), "invalid auction ID format")
	suite.mockRepo.AssertNotCalled(suite.T(), "GetAuctionById")
	suite.mockRepo.AssertNotCalled(suite.T(), "CreateBid")
}

func (suite *BidServiceTestSuite) TestCreateBid_AuctionNotFound() {
	// Arrange
	auctionID := uuid.New()
	req := &CreateBidRequest{
		Amount:    150.0,
		BidderID:  uuid.New().String(),
		AuctionID: auctionID.String(),
	}

	// Mock expectations
	suite.mockRepo.On("GetAuctionById", auctionID).Return(nil, errors.New("auction not found"))

	// Act
	result, err := suite.service.CreateBid(req)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Contains(suite.T(), err.Error(), "failed to get auction")
	suite.mockRepo.AssertCalled(suite.T(), "GetAuctionById", auctionID)
	suite.mockRepo.AssertNotCalled(suite.T(), "CreateBid")
}

func (suite *BidServiceTestSuite) TestCreateBid_InvalidBidderID() {
	// Arrange
	auctionID := uuid.New()
	req := &CreateBidRequest{
		Amount:    150.0,
		BidderID:  "invalid-uuid",
		AuctionID: auctionID.String(),
	}

	mockAuction := &database.Auction{
		ID:            auctionID,
		Title:         "Test Auction",
		Description:   "This is a test auction",
		StartingPrice: 100.0,
		SellerID:      uuid.New(),
	}

	// Mock expectations
	suite.mockRepo.On("GetAuctionById", auctionID).Return(mockAuction, nil)

	// Act
	result, err := suite.service.CreateBid(req)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Contains(suite.T(), err.Error(), "invalid bidder ID")
	suite.mockRepo.AssertCalled(suite.T(), "GetAuctionById", auctionID)
	suite.mockRepo.AssertNotCalled(suite.T(), "CreateBid")
}

func (suite *BidServiceTestSuite) TestCreateBid_SellerCannotBid() {
	// Arrange
	auctionID := uuid.New()
	sellerID := uuid.New()

	req := &CreateBidRequest{
		Amount:    150.0,
		BidderID:  sellerID.String(),
		AuctionID: auctionID.String(),
	}

	mockAuction := &database.Auction{
		ID:            auctionID,
		Title:         "Test Auction",
		Description:   "This is a test auction",
		StartingPrice: 100.0,
		SellerID:      sellerID,
	}

	// Mock expectations
	suite.mockRepo.On("GetAuctionById", auctionID).Return(mockAuction, nil)

	// Act
	result, err := suite.service.CreateBid(req)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Contains(suite.T(), err.Error(), "seller cannot bid on their own auction")
	suite.mockRepo.AssertCalled(suite.T(), "GetAuctionById", auctionID)
	suite.mockRepo.AssertNotCalled(suite.T(), "CreateBid")
}

func (suite *BidServiceTestSuite) TestCreateBid_BelowStartingPrice() {
	// Arrange
	auctionID := uuid.New()
	bidderID := uuid.New()
	sellerID := uuid.New()

	req := &CreateBidRequest{
		Amount:    80.0, // Below starting price
		BidderID:  bidderID.String(),
		AuctionID: auctionID.String(),
	}

	mockAuction := &database.Auction{
		ID:            auctionID,
		Title:         "Test Auction",
		Description:   "This is a test auction",
		StartingPrice: 100.0,
		SellerID:      sellerID,
	}

	// Mock expectations
	suite.mockRepo.On("GetAuctionById", auctionID).Return(mockAuction, nil)
	suite.mockRepo.On("GetHighestBidForAuction", auctionID).Return(nil, nil)

	// Act
	result, err := suite.service.CreateBid(req)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Contains(suite.T(), err.Error(), "bid amount must be at least 105.00 (current price 100.00 + minimum increment 5.00)")
	suite.mockRepo.AssertCalled(suite.T(), "GetAuctionById", auctionID)
	suite.mockRepo.AssertCalled(suite.T(), "GetHighestBidForAuction", auctionID)
	suite.mockRepo.AssertNotCalled(suite.T(), "CreateBid")
}

func (suite *BidServiceTestSuite) TestCreateBid_BelowMinimumIncrement() {
	// Arrange
	auctionID := uuid.New()
	bidderID := uuid.New()
	sellerID := uuid.New()

	req := &CreateBidRequest{
		Amount:    151.0, // Only $1 above current highest bid
		BidderID:  bidderID.String(),
		AuctionID: auctionID.String(),
	}

	mockAuction := &database.Auction{
		ID:            auctionID,
		Title:         "Test Auction",
		Description:   "This is a test auction",
		StartingPrice: 100.0,
		SellerID:      sellerID,
	}

	currentHighestBid := &database.Bid{
		ID:        uuid.New(),
		AuctionID: auctionID,
		BidderID:  uuid.New(),
		Amount:    150.0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Mock expectations
	suite.mockRepo.On("GetAuctionById", auctionID).Return(mockAuction, nil)
	suite.mockRepo.On("GetHighestBidForAuction", auctionID).Return(currentHighestBid, nil)

	// Act
	result, err := suite.service.CreateBid(req)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Contains(suite.T(), err.Error(), "bid amount must be at least")
	suite.mockRepo.AssertCalled(suite.T(), "GetAuctionById", auctionID)
	suite.mockRepo.AssertCalled(suite.T(), "GetHighestBidForAuction", auctionID)
	suite.mockRepo.AssertNotCalled(suite.T(), "CreateBid")
}

func (suite *BidServiceTestSuite) TestCreateBid_DatabaseError() {
	// Arrange
	auctionID := uuid.New()
	bidderID := uuid.New()
	sellerID := uuid.New()

	req := &CreateBidRequest{
		Amount:    150.0,
		BidderID:  bidderID.String(),
		AuctionID: auctionID.String(),
	}

	mockAuction := &database.Auction{
		ID:            auctionID,
		Title:         "Test Auction",
		Description:   "This is a test auction",
		StartingPrice: 100.0,
		SellerID:      sellerID,
	}

	// Mock expectations
	suite.mockRepo.On("GetAuctionById", auctionID).Return(mockAuction, nil)
	suite.mockRepo.On("GetHighestBidForAuction", auctionID).Return(nil, nil)
	suite.mockRepo.On("CreateBid", mock.Anything).Return(nil, errors.New("database error"))

	// Act
	result, err := suite.service.CreateBid(req)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Contains(suite.T(), err.Error(), "failed to create bid")
	suite.mockRepo.AssertExpectations(suite.T())
}
