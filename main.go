package main

import "log"

func main() {
	api := &API{}

	err := api.Initialize()
	if err != nil {
		panic(err)
	}

	log.Print("Server listening port :8080")

	api.Loop()
}
