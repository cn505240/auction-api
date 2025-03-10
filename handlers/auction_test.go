package handlers

import (
	"auction-api/database"
	"auction-api/service"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	mockhandlers "auction-api/mocks/auction-api/handlers"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type HandlersTestSuite struct {
	suite.Suite
	mockService *mockhandlers.AuctionService
	handlers    *Handlers
}

func (suite *HandlersTestSuite) SetupTest() {
	suite.mockService = &mockhandlers.AuctionService{
		Mock: mock.Mock{},
	}
	suite.handlers = NewHandlers(suite.mockService)
}

func TestHandlersSuite(t *testing.T) {
	suite.Run(t, new(HandlersTestSuite))
}

func (suite *HandlersTestSuite) TestCreateAuction_Success() {
	// Arrange
	sellerID := uuid.New()
	req := service.CreateAuctionRequest{
		Title:         "Test Auction",
		Description:   "This is a test auction",
		StartingPrice: 100.0,
		SellerID:      sellerID.String(),
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
	suite.mockService.On("CreateAuction", mock.MatchedBy(func(r *service.CreateAuctionRequest) bool {
		return r.Title == req.Title &&
			r.Description == req.Description &&
			r.StartingPrice == req.StartingPrice &&
			r.SellerID == req.SellerID
	})).Return(returnedAuction, nil)

	// Create request
	reqBody, _ := json.Marshal(req)
	httpReq, _ := http.NewRequest("POST", "/auctions", bytes.NewBuffer(reqBody))
	httpReq.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	// Act
	suite.handlers.CreateAuction(recorder, httpReq)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, recorder.Code)

	var responseAuction database.Auction
	err := json.Unmarshal(recorder.Body.Bytes(), &responseAuction)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), returnedAuction.ID, responseAuction.ID)
	assert.Equal(suite.T(), returnedAuction.Title, responseAuction.Title)
	assert.Equal(suite.T(), returnedAuction.Description, responseAuction.Description)
	assert.Equal(suite.T(), returnedAuction.StartingPrice, responseAuction.StartingPrice)
	assert.Equal(suite.T(), returnedAuction.SellerID, responseAuction.SellerID)

	suite.mockService.AssertExpectations(suite.T())
}

func (suite *HandlersTestSuite) TestCreateAuction_BadRequest() {
	// Arrange
	invalidJSON := `{"title": "Test", "starting_price": "not-a-number"}`

	// Create request
	httpReq, _ := http.NewRequest("POST", "/auctions", bytes.NewBufferString(invalidJSON))
	httpReq.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	// Act
	suite.handlers.CreateAuction(recorder, httpReq)

	// Assert
	assert.Equal(suite.T(), http.StatusBadRequest, recorder.Code)
	suite.mockService.AssertNotCalled(suite.T(), "CreateAuction")
}

func (suite *HandlersTestSuite) TestCreateAuction_ServiceError() {
	// Arrange
	sellerID := uuid.New()
	req := service.CreateAuctionRequest{
		Title:         "Test Auction",
		Description:   "This is a test auction",
		StartingPrice: 100.0,
		SellerID:      sellerID.String(),
	}

	// Mock expectations
	suite.mockService.On("CreateAuction", mock.Anything).Return(nil, errors.New("service error"))

	// Create request
	reqBody, _ := json.Marshal(req)
	httpReq, _ := http.NewRequest("POST", "/auctions", bytes.NewBuffer(reqBody))
	httpReq.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	// Act
	suite.handlers.CreateAuction(recorder, httpReq)

	// Assert
	assert.Equal(suite.T(), http.StatusInternalServerError, recorder.Code)
	suite.mockService.AssertExpectations(suite.T())
}
