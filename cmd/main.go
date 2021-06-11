package main

import (
	"GoStore/internal/server"
	"log"
	"runtime"
)

func main() {
	server, err := server.New("F:\\[Study]\\Projects\\Programming\\Main\\Go\\GoStore\\cmd\\config\\config.json")
	if err != nil {
		log.Fatal(err)
	}

	runtime.GOMAXPROCS(runtime.NumCPU())

	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}