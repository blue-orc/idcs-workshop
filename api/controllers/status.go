package controllers

import (
	"github.com/gorilla/mux"
	"idcs-workshop/services"
	"idcs-workshop/utilities"

	"net/http"
)

func InitStatusController(r *mux.Router) {
	r.HandleFunc("/status", statusHandler).Methods("GET")
	r.HandleFunc("/status/auth", authHandler).Methods("GET")
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	utilities.RespondOK(w)
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	a := services.Authorization{}
	ok := a.Authenticate(w, r)
	if !ok {
		return
	}
	utilities.RespondOK(w)
}
