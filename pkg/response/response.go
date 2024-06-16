package response

import (
	"encoding/json"
	"net/http"
)

type Data struct {
	Data         interface{} `json:"data"`
	Path         string      `json:"path"`
	NextPage     string      `json:"next_page"`
	PreviousPage string      `json:"previous_page"`
	Total        int         `json:"total"`
	PerPage      int         `json:"per_page"`
	CurrentPage  int         `json:"current_page"`
}

type Response struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
	Data    Data   `json:"data"`
	Version string `json:"version"`
}

type ResponseSingle struct {
	Message string      `json:"message"`
	Status  int         `json:"status"`
	Data    interface{} `json:"data"`
	Version string      `json:"version"`
}

type ResponseError struct {
	Error interface{} `json:"error"`
}

type ErrPayload struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// JSONResponse write json response format
func JSONResponse(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
