package external

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

}
