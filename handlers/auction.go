package handlers

import (
	"auction-api/service"
	"encoding/json"
	"net/http"
)

func (h *Handlers) CreateAuction(w http.ResponseWriter, r *http.Request) {
	var req service.CreateAuctionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	auction, err := h.service.CreateAuction(&req)
	// TODO: appropriate HTTP error codes
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(auction)
}
