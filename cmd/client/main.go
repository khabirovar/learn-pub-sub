package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril client...")
	rabbitConnString := "amqp://guest:guest@localhost:5672/"
	conn, err := amqp.Dial(rabbitConnString)
	if err != nil {
		log.Fatalf("Couldn't connect to rabbitmq %v", err)
	}
	defer conn.Close()
	fmt.Println("Connection was successful...")

	username, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatalf("Couldn't create welcome: %v", err)
	}

	queueName := fmt.Sprintf("%s.%s", routing.PauseKey, username)
	_, queue, err := pubsub.DeclareAndBind(conn, routing.ExchangePerilDirect, queueName, routing.PauseKey, pubsub.Transient)
	if err != nil {
		log.Fatalf("Couldn't create queue: %v", err)
	}
	fmt.Printf("Queue %v declared and bound!\n", queue.Name)

	// wait for signal ctrl + c
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	fmt.Println("Program is shuting down...")
}
