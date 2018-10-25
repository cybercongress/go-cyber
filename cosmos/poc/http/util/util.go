package util

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

func HandleError(err error, w http.ResponseWriter) {
	type ErrorJson struct {
		Error string `json:"error"`
	}
	w.WriteHeader(http.StatusInternalServerError)

	errBytes, _ := json.Marshal(ErrorJson{err.Error()})
	w.Write(errBytes)
}

func GetSingleParamValue(name string, r *http.Request) (string, error) {
	values, ok := r.URL.Query()[name]

	if !ok || values == nil {
		return "", errors.New("Cannot find param " + name)
	}

	param := ""
	if len(values) == 1 {
		param = values[0]
	}

	return param, nil
}

func GetBooleanParamValue(name string, r *http.Request) (bool, error) {
	paramStr, err := GetSingleParamValue(name, r)
	if err != nil {
		return false, nil
	}

	param, err := strconv.ParseBool(paramStr)
	if err != nil {
		return false, err
	}

	return param, nil
}
