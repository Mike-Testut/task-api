package httpjson

import(
	"encoding/json"
	"net/http"
)

type jsonResponse struct{
	Data interface{} `json:"data"`
}

type jsonError struct{
	Error string `json:"error"`
}

func WriteJSON(w http.ResponseWriter, status int,data interface{}){
	response := jsonResponse{Data:data}

	js, err := json.Marshal(response)
	if err != nil{
		http.Error(w, "internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

}

func ErrorJSON(w http.ResponseWriter, status int, message string){
	response := jsonError{Error: message}

	js, err := json.Marshal(response)
	if err != nil{
		http.Error(w, "internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)
}