package main

import (
	"net/http"
)

func main() {
	r := getRouter()
	http.ListenAndServe(":3000", r)
}
