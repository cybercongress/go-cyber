package util

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

func GetSingleParamValue(paramName string, r *http.Request) (string, error) {
	addresses, ok := r.URL.Query()[paramName]

	if !ok || addresses == nil {
		return "", errors.New("Cannot find param " + paramName)
	}

	paramValue := ""
	if len(addresses) == 1 {
		paramValue = addresses[0]
	}

	return paramValue, nil
}
