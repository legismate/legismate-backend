package main

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/legismate/legismate_backend/models"
)

func AddressCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rawAddress := r.URL.Query().Get("address")

		if rawAddress == "" {
			http.Error(w, "Missing 'address' query parameter", http.StatusBadRequest)
			return
		}

		address := models.Address{Raw: rawAddress}

		ctx := context.WithValue(r.Context(), "address", &address)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(middleware.SetHeader("Content-type", "application/json"))
	r.Use(AddressCtx)

	// everything in districts will take address and level query param
	r.Route("/districts", func(r chi.Router) { // everything in districts will take address
		r.Get("/", getDistrictByLocation)
		r.Get("/representatives", getRepsByDistrict) // query params for level adn district number
		r.Get("/deadlines", getDeadlinesByDistrict)
	})

	r.Route("/representatives", func(r chi.Router) {
		r.Get("/{repId}", getRepByID)
		r.Get("/{repId}/history", getRepBillHistory)
	})

	r.Route("/bills", func(r chi.Router) {
		r.Get("/", getBillsByLevel)
		// todo: review -- changed this to /legistar since its specific to legistar
		r.Get("/legistar/{legistarId}", getBillsByLegistarID)
	})

	return r
}
