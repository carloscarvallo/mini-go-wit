package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	wittoken := os.Getenv("WITAI_TOKEN")

	url := "https://api.wit.ai/message?v=20160526&q=cual%2520es%2520el%2520clima%2520en%2520asuncion"
	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("authorization", "Bearer "+wittoken)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	fmt.Println(res)
	fmt.Println(string(body))
}
