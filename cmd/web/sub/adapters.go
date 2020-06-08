package sub

import (
	"encoding/json"
	"net/http"
)

// response is a basic response to all requests
type response struct {
	Success   bool        `json:"success"`
	Error     string      `json:"error"`
	Data      interface{} `json:"data"`
	Status    int         `json:"-"`
	RealError string      `json:"-"`
}

// do is a helper function to do response
func (r *response) do(w http.ResponseWriter) {
	j, err := json.Marshal(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(r.Status)

	b, err := w.Write(j)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if b == 0 {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
