package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/joho/godotenv"
)

const (
	baseURL = "https://api.wit.ai/message"
)

func main() {
	// Read environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	wittoken := os.Getenv("WITAI_TOKEN")

	// Hardcoding message yet
	msg := "Cual es el clima en Asuncion?"

	// Attaching query params for the wit.ai API
	v := url.Values{}
	v.Add("v", "2016052")
	v.Add("q", msg)
	// Encode all values
	encodedValues := v.Encode()

	// Return a string with that format
	url := fmt.Sprintf("%s?%s", baseURL, encodedValues)
	// Use NewRequest for control HTTP headers in this case adding wit_token
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("authorization", "Bearer "+wittoken)
	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	// printing the body
	fmt.Println(string(body))
}
