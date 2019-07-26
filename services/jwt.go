package services

import (
	"encoding/json"
	"errors"
	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
	"idcs-workshop/config"
	"idcs-workshop/models"
	"idcs-workshop/utilities"
	"io/ioutil"
	"net/http"
	"strings"
)

type Authorization struct {
	//TODO: find a way to cache this value
	publicKey models.Key
}

type AccessTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (a *Authorization) Authenticate(w http.ResponseWriter, r *http.Request) bool {
	auth := r.Header.Get("Authorization")
	split := strings.Split(auth, " ")
	if split[0] != "Bearer" {
		utilities.RespondUnauthorized(w, "Endpoint requires Bearer token authentication")
		return false
	}
	valid, err := a.validateToken(split[1])
	if !valid || err != nil {
		utilities.RespondUnauthorized(w, "Invalid token")
		return false
	}
	return true
}

func (a *Authorization) validateToken(accessToken string) (bool, error) {
	_, err := a.GetClaims(accessToken)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (a *Authorization) GetClaims(accessToken string) (map[string]interface{}, error) {
	claims := make(map[string]interface{})
	parsedJWT, err := jwt.ParseSigned(accessToken)
	if err != nil {
		return claims, errors.New("Unable to parse access token to JWT")
	}
	jwk, err := a.getPublicKey(accessToken)
	if err != nil {
		return claims, err
	}

	err = parsedJWT.Claims(&jwk, &claims)
	if err != nil {
		return claims, errors.New("Unable to get claims/validate token " + err.Error())
	}

	return claims, nil
}

func (a *Authorization) getPublicKey(accessToken string) (jose.JSONWebKey, error) {
	publicKey, err := a.getJwtSigningKey(accessToken)
	var jwk jose.JSONWebKey
	keyBytes, err := json.Marshal(publicKey)
	if err != nil {
		return jose.JSONWebKey{}, errors.New("Cannot marshal key value to JSON: " + err.Error())
	}
	err = jwk.UnmarshalJSON(keyBytes)
	if err != nil {
		return jose.JSONWebKey{}, err
	}
	return jwk, nil
}

func (a *Authorization) getJwtSigningKey(accessToken string) (models.Key, error) {
	if a.publicKey.E != "" {
		return a.publicKey, nil
	}
	hc := http.Client{}
	req, err := http.NewRequest("GET", config.Get("IDCS_SIGNING_KEY_URL"), nil)
	if err != nil {
		return models.Key{}, errors.New("Error getting JWT signing key: " + err.Error())
	}
	req.Header.Add("Authorization", "Bearer "+accessToken)
	resp, err := hc.Do(req)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	k := models.IdcsSigningKey{}
	err = json.Unmarshal(body, &k)
	if err != nil {
		return models.Key{}, errors.New("Error deserializing IDCS JWT: " + err.Error())
	}
	if len(k.Keys) == 0 {
		return models.Key{}, errors.New("Error getting IDCS Signing Keys, no keys returned from IDCS")
	}
	a.publicKey = k.Keys[0]

	return a.publicKey, nil
}
