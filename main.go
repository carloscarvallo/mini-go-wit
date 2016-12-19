package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

var err = godotenv.Load()

var wit_token = os.Getenv("WIT_TOKEN")
var fb_token = os.Getenv("FB_TOKEN")
var app_token = os.Getenv("APP_TOKEN")

const (
	baseURL = "https://api.wit.ai"
)

// Message struct for json
type Message struct {
	ID      int    `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

func tokenVerify(w http.ResponseWriter, req *http.Request) {
	params := req.URL.Query()
	hub_mode := params.Get("hub.mode")
	verify_token := params.Get("hub.verify_token")
	fmt.Println(hub_mode, verify_token)
}

func postMessage(w http.ResponseWriter, req *http.Request) {
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
	request.Header.Add("authorization", "Bearer "+wit_token)

	res, _ := http.DefaultClient.Do(request)
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	fmt.Println(string(body))

	// write json to http.Response
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/message", postMessage).Methods("POST")
	router.HandleFunc("/webhook", tokenVerify).Methods("GET")
	log.Fatal(http.ListenAndServe(":5000", router))
}
