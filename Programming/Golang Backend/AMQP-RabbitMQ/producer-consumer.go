package main

import (
	"context"
	"fmt"
	"time"

	amq "github.com/rabbitmq/amqp091-go"
)

// exchange name + type
// queue name + routing key
// msg publish -> exchange -> (routing key) queue -> ... waiting to consume
func declareQueue(ch *amq.Channel, exchange, typee string, queue, route string) error {
	// ----	// Exchange
	if err := ch.ExchangeDeclare(exchange, typee, false, false, false, false, nil); err != nil {
		return err
	}

	// ---- // Queue
	q, err := ch.QueueDeclare(queue, true, false, false, false, nil)
	if err != nil {
		return err
	}

	fmt.Println(q.Name)

	// ---- // Bind
	if err := ch.QueueBind(q.Name, route, exchange, false, nil); err != nil {
		return err
	}
	return nil
}

func main() {

	conn, err := amq.Dial(
		"amqp://guest:guest@localhost:5672",
	)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()

	exchangeName := "my-exchange"
	exchangeType := "direct"
	routingKey := "my-key"
	queueName := "my-queue"
	if err := declareQueue(ch, exchangeName, exchangeType, queueName, routingKey); err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Publish Message
	if err := ch.PublishWithContext(ctx, exchangeName, routingKey, false, false, amq.Publishing{
		ContentType: "text/plain",
		Body:        []byte("Hello World!!!"),
	}); err != nil {
		panic(err)
	}

	fmt.Println("Waiting.... Message Published .. ..")
	time.Sleep(2 * time.Second)

	// Consume Message
	msgs, err := ch.ConsumeWithContext(ctx, queueName, "my-consumer", true, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	got := <-msgs
	fmt.Println("Consumed: ", string(got.Body))
}
