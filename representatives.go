package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/legismate/legismate_backend/cache"
	"github.com/legismate/legismate_backend/external"
	"github.com/legismate/legismate_backend/models"
	"net/http"
)

var legistar = external.GetLegistarApi("seattle")
var lCache = cache.GetLegisCache()

func getLegistarRepFromEmail(emailBase64 string) (*external.Person, error) {
	fmt.Println("retrieving legistar rep from legistar API")
	email, err := base64.StdEncoding.DecodeString(emailBase64)
	if err != nil {
		return nil, fmt.Errorf("unable to decode rep base64 email: %w", err)
	}

	person, err := legistar.GetPersonByEmail(string(email))
	if err != nil {
		return nil, fmt.Errorf("get rep from legistar using  email: %w", err)
	}
	return person, err
	//return &models.Representative{
	//	Name:   fmt.Sprintf("%s %s", person.PersonFirstName, person.PersonLastName),
	//	Party:  "Unknown",
	//	Phones: []models.Phone{models.Phone(person.PersonPhone)},
	//	Email:  person.PersonEmail,
	//	URLs:   []models.URL{models.URL(person.PersonWWW)},
	//	Office: &models.Office{LevelEnum: models.City, Name: "Seattle City Council"},
	//}, nil
}

func getRepByID(w http.ResponseWriter, r *http.Request) {
	emailBase64 := chi.URLParam(r, "repId")
	if rep, err := getLegistarRepFromEmail(emailBase64); err != nil {
		http.Error(w, "unable to get legistar rep: "+err.Error(), http.StatusInternalServerError)
	} else if err = json.NewEncoder(w).Encode(rep); err != nil {
		http.Error(w, "encoding response error: "+err.Error(), http.StatusInternalServerError)
	}
}

func getBillDetailFromVote(vote *external.Vote) (eventItem *external.EventItem, err error) {
	cacheKey := fmt.Sprintf("legistar.GetEventItemDetail:%d", vote.VoteId)
	if err = lCache.GetFromCache(cacheKey, eventItem); err != nil && !lCache.NotFound(err) {
		fmt.Printf("unexpected error hitting cache, retrieving person by legistar api\n error: %s\n", err.Error())
	}
	fmt.Println("retrieving legistar rep from legistar API")
	eventItem, err = legistar.GetEventItemDetail(vote.VoteId, vote.VoteEventItemId)
	if err != nil {
		return nil, fmt.Errorf("couldn't get vote details %w", err)
	}
	lCache.AddToCache(cacheKey, eventItem)
	return
}

// this call takes forever, probably because it has to make like fifty api calls to get the bill detail
func getRepBillHistory(w http.ResponseWriter, r *http.Request) {
	emailBase64 := chi.URLParam(r, "repId")
	rep, err := getLegistarRepFromEmail(emailBase64)
	if err != nil {
		http.Error(w, "unable to get legistar rep: "+err.Error(), http.StatusInternalServerError)
		return
	}
	votingRecord, err := legistar.GetPersonVotingRecord(rep.PersonId, 6)
	if err != nil {
		http.Error(w, fmt.Errorf("couldn't get voting record: %w", err).Error(), http.StatusInternalServerError)
		return
	}
	var representativeVotes []*models.Vote
	// add in details about the bill the vote is associated with
	for _, vote := range votingRecord {
		voteDetail, err := legistar.GetEventItemDetail(vote.VoteId, vote.VoteEventItemId)
		if err != nil {
			http.Error(w, fmt.Errorf("couldn't get vote details: %w", err).Error(), http.StatusInternalServerError)
			return
		}
		agendaDate, _ := external.ParseLegistarTime(voteDetail.EventItemLastModifiedUtc)
		representativeVotes = append(representativeVotes, &models.Vote{
			Bill: &models.Bill{
				File:       voteDetail.EventItemMatterFile,
				Title:      voteDetail.EventItemTitle,
				AgendaDate: agendaDate,
				Status:     voteDetail.EventItemMatterStatus,
				LegistarID: voteDetail.EventItemMatterId,
				//Committee:  "UNKNOWN", // todo: should we make another call for committee? idk.
			},
			RepresentativeVote: vote.VoteValueName,
		})
	}
	if err = json.NewEncoder(w).Encode(representativeVotes); err != nil {
		http.Error(w, "encoding response error: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
