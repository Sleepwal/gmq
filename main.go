package main

import (
	"gmq/server"
	"log"
)

func main() {
	config := &server.Config{
		ListenAddr: ":3000",
		StoreProducerFunc: func() server.Storage {
			return server.NewMemoryStore()
		},
	}

	s, err := server.NewServer(config)
	if err != nil {
		log.Fatal(err)
	}
	s.Start()
}
