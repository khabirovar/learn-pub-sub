package main

import (
	"fmt"
	"log"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril server...")
	rabbitConnString := "amqp://guest:guest@localhost:5672/"
	conn, err := amqp.Dial(rabbitConnString)
	if err != nil {
		log.Fatalf("Couldn't connect to rabbitmq %v\n", err)
	}
	defer conn.Close()
	fmt.Println("Connection was successful...")

	_, gameLogsCh, err := pubsub.DeclareAndBind(
		conn,
		routing.ExchangePerilTopic, 
		routing.GameLogSlug,
		fmt.Sprintf("%s.*", routing.GameLogSlug),
		pubsub.Durable,
	)
	if err != nil {
		log.Fatalf("Couldn't create queue %s: %v", routing.GameLogSlug, err)
	}
	fmt.Printf("Created queue %v\n", gameLogsCh.Name)

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Couldn't create channel to rabbitmq: %v\n", err)
	}

	gamelogic.PrintServerHelp()
	for {
		words := gamelogic.GetInput()
		if len(words) < 1 {
			continue
		}

		if words[0] == "pause" {
			fmt.Println("Sending a pause message...")
			err = pubsub.PublishJSON(ch, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{IsPaused: true})
			if err != nil {
				log.Fatalf("Publishing error: %v", err)
			}
		} else if words[0] == "resume" {
			fmt.Println("Sending a resume message...")
			err = pubsub.PublishJSON(ch, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{IsPaused: false})
			if err != nil {
				log.Fatalf("Publishing error: %v", err)
			}
		} else if words[0] == "quit" {
			fmt.Println("Exiting...")
			break
		} else {
			fmt.Println("Unknown command")
		}
	}
}
