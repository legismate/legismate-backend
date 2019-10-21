package api

import (
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"github.com/legismate/legismate_backend/cache"
)

var legisCache = cache.GetLegisCache()

func GetRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(middleware.SetHeader("Content-type", "application/json"))
	//r.Use(CacheCtx)

	// everything in districts will take address and level query param
	r.Route("/districts", func(r chi.Router) { // everything in districts will take address
		r.Get("/", getDistrictByLocation)
		r.Get("/representatives", getRepsByDistrict) // query params for level and district number
		r.Get("/deadlines", getDeadlinesByDistrict)
	})

	r.Route("/representatives", func(r chi.Router) {
		reps := getRepresentatives(legisCache)
		r.Get("/{repId}", reps.getRepByID)
		r.Get("/{repId}/history", reps.getRepBillHistory)
	})

	r.Route("/bills", func(r chi.Router) {
		r.Get("/", getBillsByLevel)
		// todo: review -- changed this to /legistar since its specific to legistar
		r.Get("/legistar/{legistarId}", getBillsByLegistarID)
	})

	return r
}
