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

type MatterVersions struct {
	Key   string
	Value string
}

type MatterText struct {
	MatterTextId              int
	MatterTextPlain           string
	MatterTextLastModifiedUtc Datetime
}

const (
	legistarBase      = "https://webapi.legistar.com/v1/%s" // %s is city/state/etc name
	matters           = legistarBase + "/matters"
	matter            = matters + "/%d"
	matterHistory     = matter + "/histories"
	matterTextVersion = matter + "/versions"
	matterText        = matter + "/texts/%d"
	person            = legistarBase + "/persons"
	personVote        = person + "/%d/votes"
)

func parseLegistarTime(lTime Datetime) (time.Time, error) {
	return time.Parse("2006-01-02T15:04:05", string(lTime))
}

func mapSingleMatterToBill(matter *Matter) (bill *models.Bill) {
	agendaDate, err := parseLegistarTime(matter.MatterAgendaDate)
	if err != nil {
		fmt.Printf("can't handle this date!! %s \n Error: %s", matter.MatterAgendaDate, err.Error())
		// todo: don't know if we should bail
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

// doSimpleAPIGetRequest will create a new request and set application/json content type header, the perform the get request.
//    if status code != ok, will return error
func doSimpleAPIGetRequest(cli *http.Client, URL string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, URL, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to create http request, err: %s", err.Error())
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := cli.Do(req)
	if err != nil {
		return nil, fmt.Errorf("unable to create get, err: %s", err.Error())
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("got bad status code " + resp.Status)
	}
	return resp, nil
}

func GetUpcomingBills(client string) ([]*models.Bill, error) {
	cli := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(matters, client), nil)
	if err != nil {
		return nil, fmt.Errorf("unable to create http request, err: %s", err.Error())
	}
	req.Header.Set("Content-Type", "application/json")
	// get everything later than today
	today := time.Now()
	q := req.URL.Query()
	q.Add("$filter", fmt.Sprintf("MatterAgendaDate ge datetime'%s'", today.Format("2006-01-02")))
	q.Add("$orderby", "MatterAgendaDate asc")
	req.URL.RawQuery = q.Encode()
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
	bills := mapMattersToBills(matters)
	return bills, nil
}

func GetSingleBillDetail(matterId int) (*models.BillDetailed, error) {
	cli := &http.Client{}
	// fixme: don't want to have to require address at every api call, but legistar does require you always pass the "Client"
	//  variable. so we have to figure something out here. for now we are hardcoding.
	resp, err := doSimpleAPIGetRequest(cli, fmt.Sprintf(matter, "seattle", matterId))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var matter Matter
	if err = json.NewDecoder(resp.Body).Decode(&matter); err != nil {
		return nil, fmt.Errorf("unable to unmarshal response! error: %s", err.Error())
	}
	fmt.Printf("matter single!! %+v", matter)
	//detailed := &models.BillDetailed{Bill: mapSingleMatterToBill(&matter)}
	return nil, nil
}
