package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Sirupsen/logrus"
)

func WriteErr(w http.ResponseWriter, err error, httpCode int) {
	logrus.Error(err.Error())

	// write error to response
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	var errMap = map[string]interface{}{
		"httpStatus": httpCode,
		"error":      err.Error(),
	}

	errJson, _ := json.Marshal(errMap)
	http.Error(w, string(errJson), httpCode)
}
