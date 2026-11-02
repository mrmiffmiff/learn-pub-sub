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
	fmt.Println("Starting Peril server...")
	conString := "amqp://guest:guest@localhost:5672/"
	con, err := amqp.Dial(conString)
	if err != nil {
		log.Fatalf("could not connect to RabbitMQ: %v", err)
	}
	defer con.Close()
	fmt.Println("Peril game server successfully connected to RabbitMQ!")

	amqpChan, err := con.Channel()
	if err != nil {
		log.Fatalf("problem creating concurrent server channel: %v", err)
	}

	gamelogic.PrintServerHelp()

	for {
		words := gamelogic.GetInput()
		if len(words) < 1 {
			continue
		}
		switch words[0] {
		case "pause":
			fmt.Println("Publishing pause state...")
			err = pubsub.PublishJSON(amqpChan, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{
				IsPaused: true,
			})
			if err != nil {
				log.Printf("could not publish time: %v", err)
			}
		case "resume":
			fmt.Println("Publishing resume state...")
			err = pubsub.PublishJSON(amqpChan, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{
				IsPaused: false,
			})
			if err != nil {
				log.Printf("could not publish time: %v", err)
			}
		case "quit":
			log.Println("Exiting...")
			return
		default:
			fmt.Println("I don't understand.")
		}
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	fmt.Println("RabbitMQ connection closed.")
	os.Exit(0)
}
