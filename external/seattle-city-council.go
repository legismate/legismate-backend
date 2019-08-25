package external

type DistrictI int

type CouncilMember struct {
	MemberName          string
	DistrictDescription string
	PhotoURL            string
	Website             string
	Email               string
	Phone               string
}

var SeattleCityCouncil = map[DistrictI]*CouncilMember{
	1: {
		"Lisa Herbold",
		"West Seattle, South Park",
		"//www.seattle.gov/Images/Council/Members/Herbold/Herbold_225x225.jpg",
		"/herbold",
		"Lisa.Herbold@seattle.gov",
		"206-684-8803",
	},
	2: {
		"Bruce Harrell",
		"South Seattle, Georgetown",
		"//www.seattle.gov/Images/Council/Members/Harrell/Bruce-Harrell-2018_225x225.jpg",
		"/harrell",
		"Bruce.Harrell@seattle.gov",
		"206-684-8804",
	},
	3: {
		"Kshama Sawant",
		"Central Seattle",
		"//www.seattle.gov/Images/Council/Members/Sawant/Sawant_225x225.jpg",
		"/sawant",
		"Kshama.Sawant@seattle.gov",
		"206-684-8016",
	},
	4: {
		"Abel Pacheco",
		"Northeast Seattle",
		"//www.seattle.gov/Images/Council/Members/Pacheco/Abel-Pacheco-2019---Bokeh-Background.jpg",
		"/pacheco",
		"Abel.Pacheco@seattle.gov",
		"206-684-8808",
	},
	5: {
		"Debora Juarez",
		"North Seattle",
		"//www.seattle.gov/Images/Council/Members/Juarez/Debora-Juarez-2018_225x225.jpg",
		"/juarez",
		"Debora.Juarez@seattle.gov",
		"206-684-8805",
	},
	6: {
		"Mike O'Brien",
		"Northwest Seattle",
		"//www.seattle.gov/Images/Council/Members/OBrien/Mike-OBrien-2018_225x225.jpg",
		"/obrien",
		"Mike.OBrien@seattle.gov",
		"206-684-8800",
	},
	7: {
		"Sally Bagshaw",
		"Pioneer Square to Magnolia",
		"//www.seattle.gov/Images/Council/Members/Bagshaw/Bagshaw_225x225.jpg",
		"/bagshaw",
		"Sally.Bagshaw@seattle.gov",
		"206-684-8801",
	},
	8: {
		"Teresa Mosqueda",
		"Citywide",
		"//www.seattle.gov/Images/Council/Members/Mosqueda/Mosqueda_225x225.jpg",
		"/mosqueda",
		"Teresa.Mosqueda@seattle.gov",
		"206-684-8806",
	},
	9: {
		"Lorena González",
		"Citywide",
		"//www.seattle.gov/Images/Council/Members/Gonzalez/Gonzalez-225x225.jpg",
		"/gonzalez",
		"Lorena.Gonzalez@seattle.gov",
		"206-684-8802",
	},
}
