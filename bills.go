package main

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"net/http"
	"strconv"

	//"github.com/go-chi/chi"

	"github.com/legismate/legismate_backend/external"
	"github.com/legismate/legismate_backend/models"
)

// getBillsByLevel will return bills by a specific level enum, if no level query parameter is passed, it will error out
//   Right now we are only doing seattle city council, so no matter what we're only using that.
func getBillsByLevel(w http.ResponseWriter, r *http.Request) {
	level := r.URL.Query().Get("level")
	levelEnum, err := models.GetLevelFromString(level)
	if err != nil || levelEnum != models.City {
		http.Error(w, "level is required parameter and must be CITY", http.StatusBadRequest)
		return
	}
	//todo: review -- levelValue is a lame query parameter
	levelValue := r.URL.Query().Get("levelValue")
	//we only take seattle right now
	if levelValue != "seattle" {
		http.Error(w, "levelValue is required query parameter and must be seattle", http.StatusBadRequest)
		return
	}
	bills, err := external.GetUpcomingBills("seattle")
	if err != nil {
		http.Error(w, "upcoming bills error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	respEncoder := json.NewEncoder(w)
	if err = respEncoder.Encode(bills); err != nil {
		http.Error(w, "encoding response error: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func getBillsByLegistarID(w http.ResponseWriter, r *http.Request) {
	matterId := chi.URLParam(r, "legistarId")
	mIdInt, err := strconv.Atoi(matterId)
	if err != nil {
		http.Error(w, "couldn't convert this to an integer: "+matterId, http.StatusBadRequest)
	}
	// fixme: don't want to have to require address at every api call, but legistar does require you always pass the "Client"
	//  variable. so we have to figure something out here. for now we are hardcoding.
	//  ** maybe everything should take address query param? **
	detailedBill, err := external.GetSingleBillDetail(mIdInt, "seattle")
	if err != nil {
		http.Error(w, "single bill detail error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if err = json.NewEncoder(w).Encode(detailedBill); err != nil {
		http.Error(w, "encoding response error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
