package p

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/ChimeraCoder/anaconda"
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

// EventRequestBody is a struct of HTTP Request Body.
type EventRequestBody struct {
	Token string `json:"token"`
	Event Event  `json:"event"`
	MinimumRequestBody
}

// Event is a struct of Slack Event
type Event struct {
	Channel string `json:"channel"`
	User    string `json:"user"`
	Text    string `json:"text"`
	MinimumRequestBody
}

// Tweet is a handler to be kicked by Cloud Function.
func Tweet(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		errorResponse(w, fmt.Sprintf("Fail to read body: %s", err))
		return
	}

	var mrb MinimumRequestBody
	if err := json.Unmarshal(body, &mrb); err != nil {
		errorResponse(w, fmt.Sprintf("Fail to decode mrb: %s", err))
		return
	}

	if mrb.Type == "url_verification" {
		var vrb VerifyRequestBody
		if err := json.Unmarshal(body, &vrb); err != nil {
			errorResponse(w, fmt.Sprintf("Fail to decode vrb: %s", err))
			return
		}
		successResponse(w, vrb.Challenge)
	} else if mrb.Type == "event_callback" {
		var erb EventRequestBody
		if err := json.Unmarshal(body, &erb); err != nil {
			errorResponse(w, fmt.Sprintf("Fail to decode erb: %s", err))
			return
		}
		if erb.Event.Channel != os.Getenv("SLACK_CHANNEL") {
			successResponse(w, "")
			return
		}
		if erb.Event.User != os.Getenv("SLACK_USER") {
			successResponse(w, "")
			return
		}
		tweet(erb.Event.Text)
		successResponse(w, "")
	} else {
		errorResponse(w, "Invalid type")
	}
	return
}

func tweet(message string) {
	if message == "" {
		return
	}
	app := anaconda.NewTwitterApiWithCredentials(
		os.Getenv("TWITTER_ACCESS_TOKEN"),
		os.Getenv("TWITTER_ACCESS_TOKEN_SECRET"),
		os.Getenv("TWITTER_CONSUMER_KEY"),
		os.Getenv("TWITTER_CONSUMER_SECRET"))
	_, err := app.PostTweet(message, url.Values{})
	if err != nil {
		log.Fatal("Failed to tweet: ", err)
	}
}

func successResponse(w http.ResponseWriter, mes string) {
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "text/plain")
	w.Write([]byte(mes))
}

func errorResponse(w http.ResponseWriter, mes string) {
	log.Fatal("Err: ", mes)
	http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
}
