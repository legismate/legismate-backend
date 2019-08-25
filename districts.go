package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/legismate/legismate_backend/external"
)

var kkc = &external.KingCountyClient{}

func getDistrictByLocation(w http.ResponseWriter, r *http.Request) {
	address := r.URL.Query().Get("address")
	districtInfo, err := kkc.GetDistrictInfoByAddress(address)

	if err != nil {
		fmt.Println(err)
	}

	json.NewEncoder(w).Encode(districtInfo)
}

func getRepsByDistrict(w http.ResponseWriter, r *http.Request) {
	address := r.URL.Query().Get("address")
	districtInfo, err := kkc.GetDistrictInfoByAddress(address)

	if err != nil {
		fmt.Println(err)
	}

	rep, found := external.SeattleCityCouncil[districtInfo.Name]

	if !found {
		fmt.Println("No rep found")
	}

	json.NewEncoder(w).Encode(rep)
}

func getDeadlinesByDistrict(w http.ResponseWriter, r *http.Request) {

}
