package main

import (
	"./app"
	"log"
	"time"
)

func task() {
	timerCh := time.Tick(time.Duration(60) * time.Second)
	for range timerCh {
		log.Print("Start task - update currencies")
		app.CurrencyParser()
		log.Print("Task - update currencies, Finished")
	}
}


func main() {
	go task()
	app.Run()

}
