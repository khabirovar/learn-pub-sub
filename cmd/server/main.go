package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril server...")
	rabbitConnString := "amqp://guest:guest@localhost:5672/"
	conn, err := amqp.Dial(rabbitConnString)
	if err != nil {
		log.Fatalf("Couldn't connect to rabbitmq %v\n",err)
	}
	defer conn.Close()
	fmt.Println("Connection was successful...")

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Couldn't create channel to rabbitmq: %v\n", err)
	}
	err = pubsub.PublishJSON(ch, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{IsPaused: true})
	if err != nil {
		log.Fatalf("Publishing error: %v", err)
	}

	// wait for signal ctrl + c
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	fmt.Println("Program is shuting down...")


}
