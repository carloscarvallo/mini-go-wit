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

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	wittoken := os.Getenv("WITAI_TOKEN")

	msg := "Cual es el clima en Asuncion?"
	v := url.Values{}
	v.Add("v", "2016052")
	v.Add("q", msg)
	encodedValues := v.Encode()
	baseURL := "https://api.wit.ai/message"

	url := fmt.Sprintf("%s?%s", baseURL, encodedValues)
	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("authorization", "Bearer "+wittoken)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	fmt.Println(res)
	fmt.Println(string(body))
}
