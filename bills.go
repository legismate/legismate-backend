package main

import (
	"encoding/json"
	"github.com/legismate/legismate_backend/models"
	"net/http"

	"github.com/legismate/legismate_backend/external"
)

// getBillsByLevel will return bills by a specific level enum, if no level query parameter is passed, it will error out
//   Right now we are only doing seattle city council, so no matter what we're only using that.
func getBillsByLevel(w http.ResponseWriter, r *http.Request) {
	respEncoder := json.NewEncoder(w)
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
	if err = respEncoder.Encode(bills); err != nil {
		http.Error(w, "encoding response error: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func getBillsByLegistarID(w http.ResponseWriter, r *http.Request) {

}
