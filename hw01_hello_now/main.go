package main

import (
	"fmt"
	"log"
	"time"

	"github.com/beevik/ntp"
)

func main() {
	currentTime := time.Now().UTC()
	timeHost := "0.beevik-ntp.pool.ntp.org"
	exactTime, err := ntp.Time(timeHost)
	if err != nil {
		log.Fatalf("failed to get exact time from host: %s", timeHost)
	}
	fmt.Println("current time:", currentTime.Round(0))
	fmt.Println("exact time:", exactTime.UTC().Round(0))
}
