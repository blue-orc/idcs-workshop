package controllers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"idcs-workshop/config"
	"idcs-workshop/models"
	"idcs-workshop/services"
	"idcs-workshop/utilities"

	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func InitUserController(r *mux.Router) {
	r.HandleFunc("/user/token", tokenHandler).Methods("POST")
	r.HandleFunc("/user/refresh", refreshTokenHandler).Methods("POST")
}

func tokenHandler(w http.ResponseWriter, r *http.Request) {
	var sir models.SignInRequest
	utilities.ReadJsonHttpBody(w, r, &sir)
	if sir.Email == "" {
		utilities.RespondBadRequest(w, "Email address required")
		return
	}

	if sir.Password == "" {
		utilities.RespondBadRequest(w, "password required")
		return
	}

	form := url.Values{}
	form.Add("grant_type", "password")
	form.Add("scope", "urn:opc:idm:__myscopes__ offline_access")
	form.Add("username", sir.Email)
	form.Add("password", sir.Password)

	hc := http.Client{}
	req, err := http.NewRequest("POST", config.Get("IDCS_TOKEN_URL"), strings.NewReader(form.Encode()))
	if err != nil {
		utilities.RespondInternalServerError(w, "Failed to create request to IDCS token endpoint: "+err.Error())
		return
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", "Basic "+config.Get("IDCS_TOKEN_URL"))
	resp, err := hc.Do(req)

	body, err := ioutil.ReadAll(resp.Body)

	a := services.Authorization{}
	atr := services.AccessTokenResponse{}
	err = json.Unmarshal(body, &atr)
	if err != nil {
		utilities.RespondUnauthorized(w, err.Error())
		return
	}
	c, err := a.GetClaims(atr.AccessToken)
	if err != nil {
		utilities.RespondUnauthorized(w, err.Error())
		return
	}

	idcsID := c["user_id"].(string)
	displayName := c["user_displayname"].(string)
	email := c["sub"].(string)
	u := models.User{}
	u.IdcsID = idcsID
	u.DisplayName = displayName
	u.Email = email
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func refreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	ct := r.Header.Get("Content-Type")
	if ct != "application/x-www-form-urlencoded" {
		utilities.RespondBadRequest(w, "Content-Type must be application/x-www-form-urlencoded")
		return
	}
	r.ParseForm()
	rt := r.Form.Get("refresh_token")
	if rt == "" {
		utilities.RespondBadRequest(w, "refresh_token is empty")
		return
	}

	form := url.Values{}
	form.Add("grant_type", "refresh_token")
	form.Add("refresh_token", rt)

	hc := http.Client{}
	req, err := http.NewRequest("POST", config.Get("IDCS_TOKEN_URL"), strings.NewReader(form.Encode()))
	if err != nil {
		utilities.RespondInternalServerError(w, "Failed to create request to IDCS token endpoint: "+err.Error())
		return
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", "Basic "+config.Get("IDCS_AUTH_SECRET"))
	resp, err := hc.Do(req)

	body, err := ioutil.ReadAll(resp.Body)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}
