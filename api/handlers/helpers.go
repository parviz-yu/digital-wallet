package handlers

import (
	"encoding/json"
	"net/http"
)

func Error(w http.ResponseWriter, r *http.Request, code int, err error) {
	if err != nil {
		Respond(w, r, code, map[string]string{"error": err.Error()})
	} else {
		Respond(w, r, code, nil)
	}

}

func Respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
