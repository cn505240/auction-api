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

type AuctionServiceTestSuite struct {
	suite.Suite
	mockRepo *mockservice.AuctionRepository
	service  *AuctionService
}

func (suite *AuctionServiceTestSuite) SetupTest() {
	suite.mockRepo = &mockservice.AuctionRepository{
		Mock: mock.Mock{},
	}
	suite.service = NewAuctionService(suite.mockRepo)
}

func TestAuctionServiceSuite(t *testing.T) {
	suite.Run(t, new(AuctionServiceTestSuite))
}

func (suite *AuctionServiceTestSuite) TestCreateAuction_Success() {
	// Arrange
	sellerID := uuid.New()
	sellerIDStr := sellerID.String()

	req := &CreateAuctionRequest{
		Title:         "Test Auction",
		Description:   "This is a test auction",
		StartingPrice: 100.0,
		SellerID:      sellerIDStr,
	}

	mockUser := &database.User{
		ID:        sellerID,
		FirstName: "Test",
		LastName:  "User",
		Email:     "test@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	expectedAuction := &database.Auction{
		Title:         req.Title,
		Description:   req.Description,
		StartingPrice: req.StartingPrice,
		SellerID:      sellerID,
	}

	returnedAuction := &database.Auction{
		ID:            uuid.New(),
		Title:         req.Title,
		Description:   req.Description,
		StartingPrice: req.StartingPrice,
		SellerID:      sellerID,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Mock expectations
	suite.mockRepo.On("GetUserById", sellerID).Return(mockUser, nil)
	suite.mockRepo.On("CreateAuction", mock.MatchedBy(func(a *database.Auction) bool {
		return a.Title == expectedAuction.Title &&
			a.Description == expectedAuction.Description &&
			a.StartingPrice == expectedAuction.StartingPrice &&
			a.SellerID == expectedAuction.SellerID
	})).Return(returnedAuction, nil)

	// Act
	result, err := suite.service.CreateAuction(req)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), returnedAuction, result)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *AuctionServiceTestSuite) TestCreateAuction_EmptyTitle() {
	// Arrange
	req := &CreateAuctionRequest{
		Title:         "",
		Description:   "This is a test auction",
		StartingPrice: 100.0,
		SellerID:      uuid.New().String(),
	}

	// Act
	result, err := suite.service.CreateAuction(req)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Contains(suite.T(), err.Error(), "title cannot be empty")
	suite.mockRepo.AssertNotCalled(suite.T(), "GetUserById")
	suite.mockRepo.AssertNotCalled(suite.T(), "CreateAuction")
}

func (suite *AuctionServiceTestSuite) TestCreateAuction_InvalidPrice() {
	// Arrange
	req := &CreateAuctionRequest{
		Title:         "Test Auction",
		Description:   "This is a test auction",
		StartingPrice: 0,
		SellerID:      uuid.New().String(),
	}

	// Act
	result, err := suite.service.CreateAuction(req)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Contains(suite.T(), err.Error(), "starting price must be greater than zero")
	suite.mockRepo.AssertNotCalled(suite.T(), "GetUserById")
	suite.mockRepo.AssertNotCalled(suite.T(), "CreateAuction")
}

func (suite *AuctionServiceTestSuite) TestCreateAuction_InvalidSellerID() {
	// Arrange
	req := &CreateAuctionRequest{
		Title:         "Test Auction",
		Description:   "This is a test auction",
		StartingPrice: 100.0,
		SellerID:      "invalid-uuid",
	}

	// Act
	result, err := suite.service.CreateAuction(req)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Contains(suite.T(), err.Error(), "invalid seller ID format")
	suite.mockRepo.AssertNotCalled(suite.T(), "GetUserById")
	suite.mockRepo.AssertNotCalled(suite.T(), "CreateAuction")
}

func (suite *AuctionServiceTestSuite) TestCreateAuction_SellerNotFound() {
	// Arrange
	sellerID := uuid.New()
	req := &CreateAuctionRequest{
		Title:         "Test Auction",
		Description:   "This is a test auction",
		StartingPrice: 100.0,
		SellerID:      sellerID.String(),
	}

	// Mock expectations
	suite.mockRepo.On("GetUserById", sellerID).Return(nil, errors.New("user not found"))

	// Act
	result, err := suite.service.CreateAuction(req)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Contains(suite.T(), err.Error(), "invalid seller ID")
	suite.mockRepo.AssertCalled(suite.T(), "GetUserById", sellerID)
	suite.mockRepo.AssertNotCalled(suite.T(), "CreateAuction")
}

func (suite *AuctionServiceTestSuite) TestCreateAuction_DatabaseError() {
	// Arrange
	sellerID := uuid.New()
	req := &CreateAuctionRequest{
		Title:         "Test Auction",
		Description:   "This is a test auction",
		StartingPrice: 100.0,
		SellerID:      sellerID.String(),
	}

	mockUser := &database.User{
		ID:        sellerID,
		FirstName: "Test",
		LastName:  "User",
		Email:     "test@example.com",
	}

	// Mock expectations
	suite.mockRepo.On("GetUserById", sellerID).Return(mockUser, nil)
	suite.mockRepo.On("CreateAuction", mock.Anything).Return(nil, errors.New("database error"))

	// Act
	result, err := suite.service.CreateAuction(req)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Contains(suite.T(), err.Error(), "failed to create auction")
	suite.mockRepo.AssertExpectations(suite.T())
}
