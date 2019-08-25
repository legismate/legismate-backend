package external

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type Location struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type CandidateGIS struct {
	Address  string   `json:"address"`
	Location Location `json:"location"`
	Score    float64  `json:"score"`
}

type GetAddressResponse struct {
	Candidates []CandidateGIS `json:"candidates"`
}

type DistrictInfo struct {
	LayerName        string `json:"layerName"`
	DisplayFieldName string `json:"displayFieldName"`
	Value            string `json:"value"`
}

type GetDistrictInfoByLocationResponse struct {
	Results []DistrictInfo
}

type KingCountyClient struct{}

var kingCountyClient = &http.Client{Timeout: 10 * time.Second}

func (k *KingCountyClient) GetGISFromAddress(address string) (getAddress GetAddressResponse, err error) {
	address = url.QueryEscape(address)
	finalURL := "https://gismaps.kingcounty.gov/arcgis/rest/services/Address/Composite_locator/GeocodeServer/findAddressCandidates?Street=" + address + "&f=json&outSR=%7B%22wkid%22%3A102100%7D"

	req, err := http.NewRequest(http.MethodGet, finalURL, nil)

	if err != nil {
		return
	}

	res, err := kingCountyClient.Do(req)

	if err != nil {
		return
	}

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return
	}

	getAddress = GetAddressResponse{}

	err = json.Unmarshal(body, &getAddress)
	if err != nil {
		return
	}

	return getAddress, nil
}

func (k *KingCountyClient) GetDistrictInfoByLocation(loc Location) (getDistrictInfo GetDistrictInfoByLocationResponse, err error) {
	finalURL := "https://gismaps.kingcounty.gov/arcgis/rest/services/Districts/KingCo_Electoral_Districts/MapServer/identify?f=json&geometry=%7B%22points%22%3A%5B%5B" + fmt.Sprintf("%.9f", loc.X) + "%2C" + fmt.Sprintf("%.9f", loc.Y) + "%5D%5D%2C%22spatialReference%22%3A%7B%22wkid%22%3A102100%7D%7D&tolerance=3&returnGeometry=false&mapExtent=%7B%22xmin%22%3A-13619782.41862936%2C%22ymin%22%3A6039789.400242523%2C%22xmax%22%3A-13615755.142701477%2C%22ymax%22%3A6042646.234174757%2C%22spatialReference%22%3A%7B%22wkid%22%3A102100%7D%7D&imageDisplay=400%2C400%2C96&geometryType=esriGeometryMultipoint&sr=102100&layers=all%3A0%2C1%2C2%2C3%2C4"

	req, err := http.NewRequest(http.MethodGet, finalURL, nil)

	if err != nil {
		return
	}

	res, err := kingCountyClient.Do(req)

	if err != nil {
		return
	}

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return
	}

	getDistrictInfo = GetDistrictInfoByLocationResponse{}

	err = json.Unmarshal(body, &getDistrictInfo)
	if err != nil {
		return
	}

	return getDistrictInfo, nil
}
