// structs for our API
package models

// Phone is a phone
type Phone string

// URL is a website URL
type URL string

type LevelEnum int

const (
	City LevelEnum = iota
	County
	State
	Federal
)

// TODO: type LevelEnum string

// Representative is a council member
type Representative struct {
	Name     string
	Party    string
	Phones   []Phone
	Email    string
	URLs     []URL
	PhotoURL URL
	Office   *Office
}

// Office is the type of office a representative holds
type Office struct {
	LevelEnum LevelEnum // CITY
	Name      string    // Seattle District 6
}

// District is a local division of government
type District struct {
	Name string
}
