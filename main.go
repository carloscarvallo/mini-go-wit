package main

import (
	"bytes"
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

var envErr = godotenv.Load()
var witToken = os.Getenv("WIT_TOKEN")
var fbToken = os.Getenv("FB_TOKEN")
var appToken = os.Getenv("APP_TOKEN")

const (
	baseURL = "https://api.wit.ai"
)

type AIResponse struct {
	Msg      string
	Entities *Entities
	Intent   string
}

type Entities struct {
	Location []struct {
		Confidence float64 `json:"confidence"`
		Type       string  `json:"type"`
		Value      string  `json:"value"`
		Suggested  bool    `json:"suggested"`
	} `json:"location,omitempty"`
	Intent []struct {
		Confidence float64 `json:"confidence"`
		Value      string  `json:"value"`
	} `json:"intent,omitempty"`
}

type Converse struct {
	Confidence float64   `json:"confidence"`
	Type       string    `json:"type"`
	Msg        string    `json:"msg"`
	Entities   *Entities `json:"entities,omitempty"`
}

// Message struct for the message ifself
type Message struct {
	Mid  string `json:"mid"`
	Seq  int    `json:"seq"`
	Text string `json:"text"`
}

// Messaging struct for more especific data about messages
type Messaging []struct {
	Sender struct {
		ID string `json:"id"`
	} `json:"sender"`
	Recipient struct {
		ID string `json:"id"`
	} `json:"recipient"`
	Timestamp int64    `json:"timestamp"`
	Message   *Message `json:"message,omitempty"`
}

// ReicevedMsg struct for the Webhook Payload
type ReicevedMsg struct {
	Object string `json:"object"`
	Entry  []struct {
		ID        string    `json:"id"`
		Time      int64     `json:"time"`
		Messaging Messaging `json:"messaging,omitempty"`
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
				if value.Message != nil {
					msgParser(messagingArr)
				} else {
					fmt.Println("webhook received unknown event")
				}
			}
		}
	}
}

func msgParser(event Messaging) {
	for _, value := range event {
		senderID := value.Sender.ID
		recipientID := value.Recipient.ID
		timeOfMessage := value.Timestamp
		message := value.Message
		fmt.Printf("\n\nReceived message for user %s and page %s at %d with Message: \n", senderID, recipientID, timeOfMessage)
		fmt.Printf("%+v", message)
		//messageID := message.Mid
		messageText := message.Text

		if senderID != "957404200975823" {
			sendToAI(senderID, messageText)
		}
	}
}

func sendToAI(senderID string, messageText string) {
	// adding uri resource
	resource := "/converse"
	u, _ := url.ParseRequestURI(baseURL)
	u.Path = resource

	// attaching query params
	v := url.Values{}
	v.Add("v", "2016052")
	v.Add("session_id", "abc321")
	v.Add("q", messageText)
	encodedValues := v.Encode()
	url := fmt.Sprintf("%s?%s", u, encodedValues)

	// make request
	request, _ := http.NewRequest("POST", url, nil)
	request.Header.Add("authorization", "Bearer "+witToken)

	// taking response
	res, _ := http.DefaultClient.Do(request)
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	fmt.Println(string(body))

	//parse JSON
	var converse Converse
	json.Unmarshal(body, &converse)

	conversationRes := converse.Msg
	conversationEnt := converse.Entities
	conversationIntArr := conversationEnt.Intent
	var conversationInt string

	for _, value := range conversationIntArr {
		conversationInt = value.Value
	}

	aiRes := AIResponse{conversationRes, conversationEnt, conversationInt}
	formatMessage(senderID, aiRes)

}

func formatMessage(senderID string, aiRes AIResponse) {
	aiResMsg := aiRes.Msg

	type Recipient struct {
		ID string `json:"id"`
	}

	type Message struct {
		Text string `json:"text"`
	}

	type msgPayload struct {
		Recipient `json:"recipient"`
		Message   `json:"message"`
	}

	data := &msgPayload{
		Recipient: Recipient{
			ID: senderID,
		},
		Message: Message{
			Text: aiResMsg,
		},
	}

	payloadBytes, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}

	postMessage(payloadBytes)
}

func postMessage(payloadBytes []byte) {
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", "https://graph.facebook.com/v2.6/me/messages?access_token="+fbToken, body)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
}

func main() {
	port := ":" + os.Args[1]

	// Handle environment vars errors
	if envErr != nil {
		log.Fatal(envErr)
	}

	router := mux.NewRouter()
	router.HandleFunc("/webhook", tokenVerify).Methods("GET")
	router.HandleFunc("/webhook", msgReceiver).Methods("POST")
	log.Fatal(http.ListenAndServe(port, router))
}
