package main

import (
	"net/http"

	"github.com/legismate/legismate_backend/external"
)

// getBillsByLevel will return bills by a specific level enum, if no level query parameter is passed, it will error out
//   Right now we are only doing seattle city council, so no matter what we're only using that.
func getBillsByLevel(w http.ResponseWriter, r *http.Request) {
	level := r.URL.Query().Get("level")
	if (level === "") {

	}
	bills, err := external.GetUpcomingBills("seattle")
}

func getBillsByLegistarID(w http.ResponseWriter, r *http.Request) {

}
