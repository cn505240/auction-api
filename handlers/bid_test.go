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

type BidHandlersTestSuite struct {
	suite.Suite
	mockService *mockhandlers.AuctionService
	handlers    *Handlers
}

func (suite *BidHandlersTestSuite) SetupTest() {
	suite.mockService = &mockhandlers.AuctionService{
		Mock: mock.Mock{},
	}
	suite.handlers = NewHandlers(suite.mockService)
}

func TestBidHandlersSuite(t *testing.T) {
	suite.Run(t, new(BidHandlersTestSuite))
}

func (suite *BidHandlersTestSuite) TestCreateBid_Success() {
	// Arrange
	bidderID := uuid.New()
	auctionID := uuid.New()
	req := service.CreateBidRequest{
		Amount:    150.0,
		BidderID:  bidderID.String(),
		AuctionID: auctionID.String(),
	}

	returnedBid := &database.Bid{
		ID:        uuid.New(),
		Amount:    req.Amount,
		BidderID:  bidderID,
		AuctionID: auctionID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Mock expectations
	suite.mockService.On("CreateBid", mock.MatchedBy(func(r *service.CreateBidRequest) bool {
		return r.Amount == req.Amount &&
			r.BidderID == req.BidderID &&
			r.AuctionID == req.AuctionID
	})).Return(returnedBid, nil)

	// Create request
	reqBody, _ := json.Marshal(req)
	httpReq, _ := http.NewRequest("POST", "/bids", bytes.NewBuffer(reqBody))
	httpReq.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	// Act
	suite.handlers.CreateBid(recorder, httpReq, auctionID.String())

	// Assert
	assert.Equal(suite.T(), http.StatusOK, recorder.Code)

	var responseBid database.Bid
	err := json.Unmarshal(recorder.Body.Bytes(), &responseBid)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), returnedBid.ID, responseBid.ID)
	assert.Equal(suite.T(), returnedBid.Amount, responseBid.Amount)
	assert.Equal(suite.T(), returnedBid.BidderID, responseBid.BidderID)
	assert.Equal(suite.T(), returnedBid.AuctionID, responseBid.AuctionID)

	suite.mockService.AssertExpectations(suite.T())
}

func (suite *BidHandlersTestSuite) TestCreateBid_BadRequest() {
	// Arrange
	invalidJSON := `{"amount": "not-a-number", "bidder_id": "123"}`
	auctionID := uuid.New().String()

	// Create request
	httpReq, _ := http.NewRequest("POST", "/bids", bytes.NewBufferString(invalidJSON))
	httpReq.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	// Act
	suite.handlers.CreateBid(recorder, httpReq, auctionID)

	// Assert
	assert.Equal(suite.T(), http.StatusBadRequest, recorder.Code)
	suite.mockService.AssertNotCalled(suite.T(), "CreateBid")
}

func (suite *BidHandlersTestSuite) TestCreateBid_ServiceError() {
	// Arrange
	bidderID := uuid.New()
	auctionID := uuid.New()
	req := service.CreateBidRequest{
		Amount:    150.0,
		BidderID:  bidderID.String(),
		AuctionID: auctionID.String(),
	}

	// Mock expectations
	suite.mockService.On("CreateBid", mock.Anything).Return(nil, errors.New("service error"))

	// Create request
	reqBody, _ := json.Marshal(req)
	httpReq, _ := http.NewRequest("POST", "/bids", bytes.NewBuffer(reqBody))
	httpReq.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	// Act
	suite.handlers.CreateBid(recorder, httpReq, auctionID.String())

	// Assert
	assert.Equal(suite.T(), http.StatusInternalServerError, recorder.Code)
	suite.mockService.AssertExpectations(suite.T())
}
