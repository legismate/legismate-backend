package legistar

import (
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/legismate/legismate_backend/models"
)

type Datetime string // todo: does json unmarshal automatically parse datetime if we set this field correctly?

func ParseLegistarTime(lTime Datetime) (time.Time, error) {
	return time.Parse("2006-01-02T15:04:05", string(lTime))
}

// Legistar's version of Representative
type Person struct {
	PersonId         int
	PersonFirstName  string
	PersonLastName   string
	PersonActiveFlag int
	PersonPhone      string
	PersonEmail      string
	PersonWWW        string
}

type Matter struct {
	MatterId              int
	MatterLastModifiedUtc Datetime
	MatterFile            string
	MatterTitle           string
	MatterTypeName        string
	MatterBodyId          int // i think this and MatterBodyName is like the committee the matter originated in?
	MatterBodyName        string
	MatterIntroDate       Datetime
	MatterAgendaDate      Datetime
	MatterStatusId        int
	MatterStatusName      string
}

type MatterVersions struct {
	Key   string
	Value string
}

type MatterText struct {
	MatterTextId              int
	MatterTextPlain           string
	MatterTextLastModifiedUtc Datetime
}

type Vote struct {
	VoteId              int // this is the event id when searching for event item detail to get matter id
	VoteGuid            string
	VoteLastModifiedUtc Datetime
	VotePersonId        int
	VotePersonName      string
	VoteValueName       string
	VoteValueId         int
	VoteResult          int
	VoteEventItemId     int
}

type EventItem struct {
	EventItemId              int
	EventItemGuid            string
	EventItemLastModifiedUtc Datetime
	EventItemEventId         int
	EventItemPassedFlag      int
	EventItemActionText      string
	EventItemMatterId        int
	EventItemTitle           string
	EventItemMatterName      string
	EventItemMatterFile      string
	EventItemMatterType      string
	EventItemMatterStatus    string
}

func mapSingleMatterToBill(matter *Matter) (bill *models.Bill) {
	agendaDate, err := ParseLegistarTime(matter.MatterAgendaDate)
	if err != nil {
		// todo: don't know if we should bail, but if we do, just return an error and log it in the caller of this function
		log.WithError(err).Errorf("can't handle this date! %s", matter.MatterAgendaDate)
	}
	bill = &models.Bill{
		File:       matter.MatterFile,
		Title:      matter.MatterTitle,
		AgendaDate: agendaDate,
		Status:     matter.MatterStatusName,
		Committee:  matter.MatterBodyName,
		LegistarID: matter.MatterId,
	}
	return
}

func mapMattersToBills(matters []*Matter) (bills []*models.Bill) {
	for _, matter := range matters {
		bills = append(bills, mapSingleMatterToBill(matter))
	}
	return bills
}
