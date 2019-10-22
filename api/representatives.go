package api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	log "github.com/sirupsen/logrus"

	"github.com/legismate/legismate_backend/cache"
	"github.com/legismate/legismate_backend/external/legistar"
	"github.com/legismate/legismate_backend/models"
)

func getRepresentatives(c *cache.LegisCache) *representatives {
	return &representatives{legistar: legistar.GetLegistarApi("seattle", c)}
}

type representatives struct {
	legistar *legistar.LegistarApi
}

func (reps *representatives) getLegistarRepFromEmail(emailBase64 string) (*legistar.Person, error) {
	log.Info("retrieving legistar rep from legistar API")
	email, err := base64.StdEncoding.DecodeString(emailBase64)
	if err != nil {
		return nil, fmt.Errorf("unable to decode rep base64 email: %w", err)
	}

	person, err := reps.legistar.GetPersonByEmail(string(email))
	if err != nil {
		return nil, fmt.Errorf("get rep from legistar using  email: %w", err)
	}
	return person, err
}

func (reps *representatives) getRepByID(w http.ResponseWriter, r *http.Request) {
	emailBase64 := chi.URLParam(r, "repId")
	if rep, err := reps.getLegistarRepFromEmail(emailBase64); err != nil {
		http.Error(w, "unable to get legistar rep: "+err.Error(), http.StatusInternalServerError)
	} else if err = json.NewEncoder(w).Encode(rep); err != nil {
		http.Error(w, "encoding response error: "+err.Error(), http.StatusInternalServerError)
	}
}

// this call takes forever, probably because it has to make like fifty api calls to get the bill detail
func (reps *representatives) getRepBillHistory(w http.ResponseWriter, r *http.Request) {
	emailBase64 := chi.URLParam(r, "repId")
	rep, err := reps.getLegistarRepFromEmail(emailBase64)
	if err != nil {
		http.Error(w, "unable to get legistar rep: "+err.Error(), http.StatusInternalServerError)
		return
	}
	votingRecord, err := reps.legistar.GetPersonVotingRecord(rep.PersonId, 6)
	if err != nil {
		http.Error(w, fmt.Errorf("couldn't get voting record: %w", err).Error(), http.StatusInternalServerError)
		return
	}
	var representativeVotes []*models.Vote
	// add in details about the bill the vote is associated with
	for _, vote := range votingRecord {
		voteDetail, err := reps.legistar.GetEventItemDetail(vote.VoteId, vote.VoteEventItemId)
		if err != nil {
			http.Error(w, fmt.Errorf("couldn't get vote details: %w", err).Error(), http.StatusInternalServerError)
			return
		}
		agendaDate, _ := legistar.ParseLegistarTime(voteDetail.EventItemLastModifiedUtc)
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
