package pubsub

import (
	"context"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type SimpleQueueType int
const (
	Durable SimpleQueueType = iota
	Transient 
)

func PublishJSON[T any](ch *amqp.Channel, exchange, key string, val T) error {
	data, err := json.Marshal(val)
	if err != nil {
		return fmt.Errorf("Error creating json: %v", err)
	}

	msg := amqp.Publishing{
		ContentType: "application/json",
		Body: data,
	}
	err = ch.PublishWithContext(context.Background(),exchange, key, false, false, msg)
	if err != nil {
		return fmt.Errorf("Error publish message: %v", err)
	}
	return nil
}

func DeclareAndBind(
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType, // SimpleQueueType is an "enum" type I made to represent "durable" or "transient"
) (*amqp.Channel, amqp.Queue, error) {
	ch, err := conn.Channel()
	if err != nil {
		return  nil, amqp.Queue{}, err
	}

	var durable, autoDelete, exclusive, noWait bool
	if queueType == Durable {
		 durable = true 
	} else if queueType == Transient {
		autoDelete = true
		exclusive = true
	}
	queue, err := ch.QueueDeclare(queueName, durable, autoDelete, exclusive, noWait, nil)
	if err != nil {
		return nil, amqp.Queue{}, err
	}

	err = ch.QueueBind(queueName, key, exchange, noWait, nil)
	if err != nil {
		return nil, amqp.Queue{}, err
	}

	return ch, queue, nil
}