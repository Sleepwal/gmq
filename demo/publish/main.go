package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
)

func main() {
	url := "http://localhost:3000/publish"
	topics := []string{"topic1", "topic2"}

	for i := 0; i < 1000; i++ {
		topic := topics[rand.Intn(len(topics))]
		payload := fmt.Sprintf("data_%d", i)
		resp, err := http.Post(url+"/"+topic+"/"+payload, "application/octet-stream", nil)
		if err != nil {
			log.Fatal(err)
		}
		if resp.StatusCode != http.StatusOK {
			log.Fatal("http status error: ", resp.Status)
		}
	}
}
