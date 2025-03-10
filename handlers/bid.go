package handlers

import (
	"auction-api/service"
	"encoding/json"
	"net/http"
)

func (h *Handlers) CreateBid(w http.ResponseWriter, r *http.Request, auctionID string) {
	req := service.CreateBidRequest{
		AuctionID: auctionID,
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	bid, err := h.service.CreateBid(&req)
	// TODO: appropriate HTTP error codes
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bid)
}
