package external

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/legismate/legismate_backend/cache"
	"github.com/legismate/legismate_backend/models"
)

var lCache = cache.GetLegisCache()

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

const (
	legistarBase       = "https://webapi.legistar.com/v1/%s" // %s is city/state/etc name
	matters            = legistarBase + "/matters"
	matter             = matters + "/%d"
	matterHistory      = matter + "/histories"
	matterTextVersions = matter + "/versions"
	matterText         = matter + "/texts/%s"
	person             = legistarBase + "/persons"
	personVote         = person + "/%d/votes"
	events             = legistarBase + "/Events"
	eventItems         = events + "/%d/EventItems/%d"
)

func ParseLegistarTime(lTime Datetime) (time.Time, error) {
	return time.Parse("2006-01-02T15:04:05", string(lTime))
}

func mapSingleMatterToBill(matter *Matter) (bill *models.Bill) {
	agendaDate, err := ParseLegistarTime(matter.MatterAgendaDate)
	if err != nil {
		log.WithError(err).Errorf("can't handle this date! %s", matter.MatterAgendaDate)
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

func GetLegistarApi(city string, cache *cache.LegisCache) *LegistarApi {
	return &LegistarApi{client: city, cache: cache}
}

// LegistarApi provides methods for retrieving reps & bill events from legistar (http://webapi.legistar.com)
type LegistarApi struct {
	// client is the city/county/state that we are requesting from. we may not want to include this... i'm not sure yet
	client string
	cache  *cache.LegisCache
}

func (l *LegistarApi) formatUrl(pathFmt string, args ...interface{}) string {
	args = append([]interface{}{l.client}, args...)
	return fmt.Sprintf(pathFmt, args...)
}

// GetUpcomingBills will return a slice of bills that have an agenda date on or after today
func (l *LegistarApi) GetUpcomingBills() ([]*models.Bill, error) {
	cli := &http.Client{}
	today := time.Now()
	req, err := http.NewRequest(http.MethodGet, l.formatUrl(matters), nil)
	if err != nil {
		return nil, fmt.Errorf("unable to create http request, err: %s", err.Error())
	}
	req.Header.Set("Content-Type", "application/json")
	// get everything later than today, and order the results so that the most recent items are first in the slice
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

// GetSingleBillDetail will grab all the standard bill information that comes from GetUpcomingBills, but it will
// 	Also return the most recently updated time as well as the full text of the bill, and what version it is on
func (l *LegistarApi) GetSingleBillDetail(matterId int) (*models.BillDetailed, error) {
	cli := &http.Client{}
	resp, err := doSimpleAPIGetRequest(cli, l.formatUrl(matter, matterId))
	if err != nil {
		return nil, err
	}
	// ensure that the body is closed. we close it after every successful marshaling of the response into a struct,
	// this is for the unhappy paths
	defer func() {
		if resp != nil {
			resp.Body.Close()
		}
	}()
	var matter Matter
	if err = json.NewDecoder(resp.Body).Decode(&matter); err != nil {
		return nil, fmt.Errorf("unable to unmarshal response! error: %s", err.Error())
	}
	resp.Body.Close()
	// create filled out bill detail model and add extra data from matters api call
	detailed := &models.BillDetailed{Bill: mapSingleMatterToBill(&matter)}
	introDate, err := ParseLegistarTime(matter.MatterIntroDate)
	if err != nil {
		log.WithError(err).Error("unable to parse time!")
	}
	detailed.IntroducedDate = introDate

	// get latest version id of bill to get text body and set current version number
	resp, err = doSimpleAPIGetRequest(cli, l.formatUrl(matterTextVersions, matterId))
	if err != nil {
		return nil, err
	}
	var versions []*MatterVersions
	if err = json.NewDecoder(resp.Body).Decode(&versions); err != nil {
		return nil, fmt.Errorf("unable to decode versions response! error: %s", err.Error())
	}
	resp.Body.Close()

	// grab the last version item, its ordered by "Value" field so the last will always be the most current
	latestVersion := versions[len(versions)-1]
	detailed.CurrentVersionNumber = latestVersion.Value
	resp, err = doSimpleAPIGetRequest(cli, l.formatUrl(matterText, matterId, latestVersion.Key))
	if err != nil {
		return nil, err
	}
	var body MatterText
	if err = json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, fmt.Errorf("unable to decode bill text response! error: %s", err.Error())
	}
	resp.Body.Close()

	detailed.FullText = body.MatterTextPlain
	return detailed, nil
}

func (l *LegistarApi) GetPersonByEmail(email string) (*Person, error) {
	legistarPerson := &Person{}
	cacheKey := fmt.Sprintf("GetPersonByEmail:%s", email)
	if err := lCache.GetFromCache(cacheKey, legistarPerson); err != nil && !lCache.NotFound(err) {
		log.WithError(err).Error("unexpected error hitting cache, retrieving person by legistar api")
	} else if err == nil {
		return legistarPerson, err
	}
	cli := &http.Client{}
	filter := fmt.Sprintf("PersonEmail eq '%s'", email)
	escapedFilter := url.QueryEscape(filter)
	resp, err := doSimpleAPIGetRequest(cli, l.formatUrl(person)+"?$filter="+escapedFilter)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var persons []Person
	if err = json.NewDecoder(resp.Body).Decode(&persons); err != nil {
		return nil, fmt.Errorf("unable to unmarshal response! error: %s", err.Error())
	}

	if len(persons) == 0 {
		return nil, fmt.Errorf("no person by e-mail %s", email)
	}
	legistarPerson = &persons[0]
	lCache.AddToCache(cacheKey, legistarPerson)
	return legistarPerson, nil
}

// GetPersonVotingRecord will return the last 50 records of votes for a rep for last number of months
func (l *LegistarApi) GetPersonVotingRecord(personId int, months int) (votes []*Vote, err error) {
	cli := &http.Client{}
	sixMonthsAgo := time.Now().AddDate(0, -months, 0)

	req, err := http.NewRequest(http.MethodGet, l.formatUrl(personVote, personId), nil)
	if err != nil {
		return nil, fmt.Errorf("unable to create http request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	// get everything later than today, and order the results so that the most recent items are first in the slice
	q := req.URL.Query()
	q.Add("$filter", fmt.Sprintf("VoteLastModifiedUtc ge datetime'%s'", sixMonthsAgo.Format("2006-01-02")))
	q.Add("$orderby", "VoteLastModifiedUtc desc")
	//todo: function param for top?
	q.Add("$top", "50")
	req.URL.RawQuery = q.Encode()
	println("legistar api string: " + req.URL.String())
	resp, err := cli.Do(req)
	if err != nil {
		return nil, fmt.Errorf("unable to get voting history: %w", err)
	}
	defer resp.Body.Close()
	if err = json.NewDecoder(resp.Body).Decode(&votes); err != nil {
		return nil, fmt.Errorf("unable to unmarshal response! %w", err)
	}
	return
}

func (l *LegistarApi) GetEventItemDetail(eventId, eventItemId int) (*EventItem, error) {
	cli := &http.Client{}
	eventItem := &EventItem{}
	cacheKey := fmt.Sprintf("GetEventItemDetail:%d", eventItemId)
	if err := lCache.GetFromCache(cacheKey, eventItem); err != nil && !lCache.NotFound(err) {
		log.WithError(err).Error("unexpected error hitting cache, retrieving event item detail by legistar api")
	} else if err == nil {
		return eventItem, err
	}
	resp, err := doSimpleAPIGetRequest(cli, l.formatUrl(eventItems, eventId, eventItemId))
	if err != nil {
		return nil, fmt.Errorf("unable to get event item detail: %w", err)
	}
	defer resp.Body.Close()
	if err = json.NewDecoder(resp.Body).Decode(eventItem); err != nil {
		return nil, fmt.Errorf("unable to unmarshal response! %w", err)
	}
	lCache.AddToCache(cacheKey, eventItem)
	return eventItem, nil
}
