package main

import (
	"github.com/legismate/legismate_backend/api"
	"net/http"
)

func main() {
	r := api.GetRouter()
	http.ListenAndServe(":3000", r)
}
