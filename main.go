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

const (
	baseURL = "https://api.wit.ai"
)

var err = godotenv.Load()
var wittoken = os.Getenv("WITAI_TOKEN")

// Message struct for json
type Message struct {
	ID      int    `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

func postMessage(w http.ResponseWriter, req *http.Request) {
	var msg Message
	dec := json.NewDecoder(req.Body)
	decErr := dec.Decode(&msg)
	if decErr != nil {
		log.Fatal(decErr)
	}
	// Attaching query params
	resource := "/message"
	u, _ := url.ParseRequestURI(baseURL)
	u.Path = resource
	v := url.Values{}
	v.Add("v", "2016052")
	v.Add("q", msg.Message)

	encodedValues := v.Encode()
	url := fmt.Sprintf("%s?%s", u, encodedValues)
	fmt.Println("url", url)
	request, _ := http.NewRequest("GET", url, nil)
	request.Header.Add("authorization", "Bearer "+wittoken)
	res, _ := http.DefaultClient.Do(request)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/message", postMessage).Methods("POST")
	log.Fatal(http.ListenAndServe(":12345", router))
}
