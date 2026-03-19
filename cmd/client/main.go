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

	gameState := gamelogic.NewGameState(username)
	queueName := fmt.Sprintf("%s.%s", routing.PauseKey, username)
	// _, queue, err := pubsub.DeclareAndBind(conn, routing.ExchangePerilDirect, queueName, routing.PauseKey, pubsub.Transient)
	err = pubsub.SubscribeJSON(conn, routing.ExchangePerilDirect, queueName, routing.PauseKey, pubsub.Transient, handlerPause(gameState))
	if err != nil {
		log.Fatalf("Couldn't subscribe queue: %v", err)
	}

	for {
		words := gamelogic.GetInput()
		if len(words) < 1 {
			continue
		}

		if words[0] == "spawn" {
			err = gameState.CommandSpawn(words)
			if err != nil {
				fmt.Printf("Couldn't spawn: %v", err)
			}
		} else if words[0] == "move" {
			_, err := gameState.CommandMove(words)
			if err != nil {
				fmt.Printf("Couldn't move: %v", err)
			}
		} else if words[0] == "status" {
			gameState.CommandStatus()
		} else if words[0] == "help" {
			gamelogic.PrintClientHelp()
		} else if words[0] == "spam" {
			fmt.Println("Spamming not allowed yet!")
		} else if words[0] == "quit" {
			gamelogic.PrintQuit()
			break
		} else {
			fmt.Println("Unknown command")
		}
	}
}

func handlerPause(gs *gamelogic.GameState) func(routing.PlayingState) {
	return func(ps routing.PlayingState) {
		defer fmt.Print("> ")
		gs.HandlePause(ps)
	}
}