package main

import (
	"Gameserver-Query"
	"log"
)

func main() {
	resp, err := gquery.Query("minecraft", "localhost", 25565)

	if err != nil {
		panic(err)
	}

	log.Println(resp)
}
