package p

import (
	"encoding/json"
	"log"
	"net/http"
)

// MinimumRequestBody is a struct of HTTP Request Body.
type MinimumRequestBody struct {
	Type string `json:"type"`
}

// VerifyRequestBody is a struct of HTTP Request Body.
type VerifyRequestBody struct {
	Token     string `json:"token"`
	Challenge string `json:"challenge"`
	MinimumRequestBody
}

// Tweet is a handler to be kicked by Cloud Function.
func Tweet(w http.ResponseWriter, r *http.Request) {
	var body = json.NewDecoder(r.Body)
	var erb interface{}
	if err := body.Decode(&erb); err != nil {
		errorResponse(w, "aaaa")
		return
	}
	log.Println(erb)
	/*
		var mrb MinimumRequestBody
		if err := body.Decode(&mrb); err != nil {
			errorResponse(w, fmt.Sprintf("Fail to decode mrb: %s", err))
			return
		}
		if mrb.Type == "url_verification" {
			var vrb VerifyRequestBody
			if err := body.Decode(&vrb); err != nil {
				errorResponse(w, fmt.Sprintf("Fail to decode vrb: %s", err))
				return
			}
			successResponse(w, vrb.Challenge)
		} else if mrb.Type == "event_callback" {
			var erb interface{}
			if err := body.Decode(&erb); err != nil {
				errorResponse(w, fmt.Sprintf("Fail to decode erb: %s", err))
				return
			}
			log.Println("Request: ", erb)
			successResponse(w, "")
		} else {
			errorResponse(w, "Invalid type")
		}
	*/
	return
}

func successResponse(w http.ResponseWriter, mes string) {
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "text/plain")
	w.Write([]byte(mes))
}

func errorResponse(w http.ResponseWriter, mes string) {
	log.Println("Err: ", mes)
	http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
}
