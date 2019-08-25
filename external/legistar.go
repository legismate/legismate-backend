package external

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/legismate/legismate_backend/models"
)

type Datetime string // todo: does json unmarshal automatically parse datetime if we set this field correctly?

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

const (
	legistarBase  = "https://webapi.legistar.com/v1/%s" // %s is city/state/etc name
	matters       = legistarBase + "/matters"
	matterHistory = matters + "/%d/histories"
	person        = legistarBase + "/persons"
	personVote    = person + "/%d/votes"
)

func mapMattersToBills(matters []*Matter) (bills []*models.Bill) {
	for _, matter := range matters {
		bills = append(bills, &models.Bill{
			File:       matter.MatterFile,
			Title:      matter.MatterTitle,
			AgendaDate: nil, // todo: add in parsing
			Status:     matter.MatterStatusName,
			Committee:  matter.MatterBodyName,
		})
	}
	return bills
}

func GetUpcomingBills(client string) ([]*models.Bill, error) {
	cli := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(matters, client), nil)
	if err != nil {
		return nil, fmt.Errorf("unable to create http request, err: %s", err.Error())
	}
	req.Header.Set("Content-Type", "application/json")
	// get everything later than today
	// MatterAgendaDate+ge+datetime%272019-08-01%27
	today := time.Now()
	q := req.URL.Query()
	q.Add("$filter", fmt.Sprintf("MatterAgendaDate ge datetime'%s'", today.Format("2006-01-02")))
	req.URL.RawQuery = q.Encode()
	fmt.Println()
	resp, err := cli.Do(req)
	if err != nil {
		return nil, fmt.Errorf("unable to create get, err: %s", err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("got bad status code " + resp.Status)
	}
	var matters []*Matter
	if err = json.NewDecoder(resp.Body).Decode(&matters); err != nil {
		return nil, fmt.Errorf("unable to decode response! error: %s", err.Error())
	}
	fmt.Println("matters: %+v", matters)
	bills := mapMattersToBills(matters)
	return bills, nil
}
