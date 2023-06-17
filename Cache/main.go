package main

import (
	"Cache/server"
	"log"
)

func main() {
	run := server.RunServer()
	if run != nil {
		log.Println(run)
	}
}
