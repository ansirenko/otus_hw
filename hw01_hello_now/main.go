package main

import (
	"fmt"
	"log"
	"time"

	"github.com/beevik/ntp"
)

func main() {
	currentTime := time.Now().Round(0).UTC()
	exactTime, err := ntp.Time("0.beevik-ntp.pool.ntp.org")
	if err != nil {
		log.Fatalf("failed to get exact time")
	}
	fmt.Println("current time:", currentTime)
	fmt.Println("exact time:", exactTime.UTC().Round(0))
}
