package external

import (
	"fmt"
	"net/http"
)

type Location struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type CandidateGIS struct {
	Address  string   `json:"address"`
	Location Location `json:"location"`
	Score    int      `json:"score"`
}

type GetAddressResponse struct {
	Candidates []CandidateGIS `json:"candidates"`
}

type KingCountyClient struct{}

func (k *KingCountyClient) getGISFromAddress(address string) Location {
	url := fmt.Sprintf("https://gismaps.kingcounty.gov/arcgis/rest/services/Address/Composite_locator/GeocodeServer/findAddressCandidates?Street=%s&f=json&outSR=%7B%22wkid%22%3A102100%7D", address)
	resp, err := http.Get(url)

}
