package core

import (
	"encoding/json"
	"errors"
	"net/http"
)

func HandleError(err error, w http.ResponseWriter) {
	type ErrorJson struct {
		Error string `json:"error"`
	}
	w.WriteHeader(http.StatusInternalServerError)

	errBytes, _ := json.Marshal(ErrorJson{err.Error()})
	w.Write(errBytes)
}

func getSingleParamValue(paramName string, r *http.Request) (string, error) {
	addresses, ok := r.URL.Query()[paramName]

	if !ok || len(addresses[0]) < 1 {
		return "", errors.New("Cannot find param " + paramName)
	}

	return addresses[0], nil
}
