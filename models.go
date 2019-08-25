package main

// Phone is a phone
type Phone string

// URL is a website URL
type URL string

// TODO: type LevelEnum string

// Representative is a council member
type Representative struct {
	Name     string
	Party    string
	Phones   []Phone
	Email    string
	URLs     []URL
	PhotoURL URL
	Office   Office
}

// Office is the type of office a representative holds
type Office struct {
	LevelEnum string // CITY
	Name      string // Seattle District 6
}

// District is a local division of government
type District struct {
	Name string
}
