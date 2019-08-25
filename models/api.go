// structs for our API
package models

import (
	"fmt"
	"time"
)

// Phone is a phone
type Phone string

// URL is a website URL
type URL string

type LevelEnum int

func (l LevelEnum) String() string {
	var levelMap = map[LevelEnum]string{
		City:    "CITY",
		County:  "COUNTY",
		State:   "STATE",
		Federal: "FEDERAL",
	}
	return levelMap[l]
}

func GetLevelFromString(levelString string) (LevelEnum, error) {
	var levelMap = map[string]LevelEnum{
		"CITY":    City,
		"COUNTY":  County,
		"STATE":   State,
		"FEDERAL": Federal,
	}
	enum, ok := levelMap[levelString]
	if !ok {
		return 0, fmt.Errorf("%s is not a string representation of the level enum", levelString)
	}
	return enum, nil
}

const (
	City LevelEnum = iota
	County
	State
	Federal
)

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

type Bill struct {
	File       string // eg Res 31894
	Title      string // eg A RESOLUTION relating to the funding of priority projects in the 2019-2024 Bicycle Master Plan Implementation Plan; requesting that the Mayor commit to building out the Bicycle Master Plan and identify funding for priority Bicycle Master Plan projects in the Mayor’s 2020 Proposed Budget.
	AgendaDate time.Time
	Status     string // could be an enum when we figure out what the valuable ones are. currently corresponds to MatterStatusName
	Committee  string // eg Sustainability and Transportation Committee
	LegistarID int    // from legistar
}
