package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

var err = godotenv.Load()

var witToken = os.Getenv("WIT_TOKEN")
var fbToken = os.Getenv("FB_TOKEN")
var appToken = os.Getenv("APP_TOKEN")

const (
	baseURL = "https://api.wit.ai"
)

// Message struct for the nested message coming from the json req
type Message struct {
	Mid  string `json:"mid"`
	Seq  int    `json:"seq"`
	Text string `json:"text"`
}

// ReicevedMsg struct from the Webhook
type ReicevedMsg struct {
	Object string `json:"object"`
	Entry  []struct {
		ID        string `json:"id"`
		Time      int64  `json:"time"`
		Messaging []struct {
			Sender struct {
				ID string `json:"id"`
			} `json:"sender"`
			Recipient struct {
				ID string `json:"id"`
			} `json:"recipient"`
			Timestamp int64   `json:"timestamp"`
			Message   Message `json:"message,omitempty"`
		} `json:"messaging,omitempty"`
	} `json:"entry"`
}

func tokenVerify(w http.ResponseWriter, req *http.Request) {
	params := req.URL.Query()
	hubMode := params.Get("hub.mode")
	verifyToken := params.Get("hub.verify_token")
	challenge := params.Get("hub.challenge")

	if hubMode == "subscribe" && verifyToken == appToken {
		fmt.Println("validating Webhook")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(challenge))
	} else {
		fmt.Println("Failed Validation. Make sure the token match")
	}
}

func msgReceiver(w http.ResponseWriter, req *http.Request) {
	var msg ReicevedMsg
	dec := json.NewDecoder(req.Body)
	decErr := dec.Decode(&msg)
	if decErr != nil {
		log.Fatal(decErr)
	}
	if msg.Object == "page" {
		entryArr := msg.Entry
		for _, value := range entryArr {
			//fmt.Println(value.ID)
			//fmt.Println(value.Time)
			messagingArr := value.Messaging
			for _, value := range messagingArr {
				if (Message{}) != value.Message {
					fmt.Println("message is came")
				} else {
					fmt.Println("webhook received unknown event")
				}
			}
		}
	}
}

func postMessage(w http.ResponseWriter, req *http.Request) {
	/*
		var msg Message
		dec := json.NewDecoder(req.Body)
		decErr := dec.Decode(&msg)
		if decErr != nil {
			log.Fatal(decErr)
		}

		// adding uri resource
		resource := "/message"
		u, _ := url.ParseRequestURI(baseURL)
		u.Path = resource

		// attaching query params
		v := url.Values{}
		v.Add("v", "2016052")
		v.Add("q", msg.Message)
		encodedValues := v.Encode()
		url := fmt.Sprintf("%s?%s", u, encodedValues)

		// make request
		request, _ := http.NewRequest("GET", url, nil)
		request.Header.Add("authorization", "Bearer "+witToken)

		res, _ := http.DefaultClient.Do(request)
		defer res.Body.Close()
		body, _ := ioutil.ReadAll(res.Body)

		fmt.Println(string(body))

		// write json to http.Response
		w.Header().Set("Content-Type", "application/json")
		w.Write(body)
	*/
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/message", postMessage).Methods("POST")
	router.HandleFunc("/webhook", tokenVerify).Methods("GET")
	router.HandleFunc("/webhook", msgReceiver).Methods("POST")
	log.Fatal(http.ListenAndServe(":5000", router))
}
