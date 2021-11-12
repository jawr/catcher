package http

import (
	"encoding/json"
	"net/http"
)

func writeJSON(response http.ResponseWriter, o interface{}) error {
	data, err := json.Marshal(o)
	if err != nil {
		return err
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
	response.Write(data)
	
	return nil
}
