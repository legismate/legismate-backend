package main

import (
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func getRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/districts", func(r chi.Router) { // everything in districts will take address
		r.Get("/", getDistrictByLocation)
		r.Get("/representatives", getRepsByDistrict) // query params for level adn district number
		r.Get("/deadlines", getDeadlinesByDistrict)
	})

	r.Route("/representatives", func(r chi.Router) {
		r.Get("/{repId}", getRepById)
		r.Get("/{repId}/history", getRepBillHistory)
	})

	r.Route("/bills", func(r chi.Router) {
		r.Get("/", getBillsByLevel)
		r.Get("/{legistarId}", getBillsByLegistarId)
	})

	return r
}
