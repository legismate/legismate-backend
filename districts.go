package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/legismate/legismate_backend/external"
	"github.com/legismate/legismate_backend/models"
)

var kkc = &external.KingCountyClient{}

func getDistrictByLocation(w http.ResponseWriter, r *http.Request) {
	address := r.URL.Query().Get("address")
	getAddressResponse, err := kkc.GetGISFromAddress(address)

	if err != nil {
		fmt.Println(err)
	}

	if len(getAddressResponse.Candidates) == 0 {
		fmt.Println("No results found")
	}

	loc := getAddressResponse.Candidates[0].Location

	getDistrictInfoResponse, err := kkc.GetDistrictInfoByLocation(loc)

	if err != nil {
		fmt.Println(err)
	}

	if len(getDistrictInfoResponse.Results) == 0 {
		fmt.Println("No results found")
	}

	districtInfo := models.District{
		Name: getDistrictInfoResponse.Results[len(getDistrictInfoResponse.Results)-1].Value,
	}

	json.NewEncoder(w).Encode(districtInfo)
}

func getRepsByDistrict(w http.ResponseWriter, r *http.Request) {

}

func getDeadlinesByDistrict(w http.ResponseWriter, r *http.Request) {

}
