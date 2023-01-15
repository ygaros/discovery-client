package client

import (
	"log"
	"strings"
	"time"
)

const duplicate = "duplicate"

func RegisterLoop(client Client) {
	ticker := time.NewTicker(5 * time.Second)
	for range ticker.C {
		err := client.Register()
		if err == nil {
			log.Println("Registered successfully in discovery server!")
			break
		} else if strings.Contains(err.Error(), duplicate) {
			log.Println("Service already registered!")
			break
		} else {
			log.Println(err)
		}
	}
}

func DoHeartBeat(client Client) {
	log.Println("Heartbeat started at 2/1min rate")
	for range time.Tick(30 * time.Second) {
		err := client.HeartBeat()
		if err != nil {
			log.Println(err)
			RegisterLoop(client)
		}
	}
}
