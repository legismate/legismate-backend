package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/legismate/legismate_backend/external"
)

func getRepByID(w http.ResponseWriter, r *http.Request) {
	emailBase64 := chi.URLParam(r, "repId")
	fmt.Printf("base64: %s", emailBase64)
	email, err := (base64.StdEncoding.DecodeString(emailBase64))

	if err != nil {
		http.Error(w, "decoding rep id error: "+err.Error(), http.StatusBadRequest)
	}

	person, err := external.GetPersonByEmail(string(email))

	if err != nil {
		http.Error(w, "find person error", http.StatusInternalServerError)
		return
	}

	if err = json.NewEncoder(w).Encode(person); err != nil {
		http.Error(w, "encoding response error: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func getRepBillHistory(w http.ResponseWriter, r *http.Request) {

}
