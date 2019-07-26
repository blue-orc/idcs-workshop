package utilities

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func ReadJsonHttpBody(w http.ResponseWriter, r *http.Request, i interface{}) error {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read http request body: "+err.Error(), http.StatusBadRequest)
	}

	err = json.Unmarshal(body, i)
	if err != nil {
		http.Error(w, "Unable to deserialize request body: "+err.Error(), http.StatusBadRequest)
	}
	return nil
}
