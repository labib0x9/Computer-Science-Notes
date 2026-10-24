package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	amq "github.com/rabbitmq/amqp091-go"
)

// exchange name + type
// queue name + routing key
// msg publish -> exchange -> (routing key) queue -> ... waiting to consume
func declareQueue(conn *amq.Connection, exchange string, queue, route string) error {
	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	// _, err = ch.QueueDelete(queue, false, false, false)
	// var amqpErr *amq.Error
	// if !(errors.As(err, &amqpErr) && amqpErr.Code == amq.NotFound) {
	// 	return err
	// }

	q, err := ch.QueueDeclare(queue, true, false, false, false, nil)
	if err != nil {
		return err
	}

	fmt.Println(q.Name)

	if err := ch.QueueBind(q.Name, route, exchange, false, nil); err != nil {
		return err
	}
	return nil
}

// remove exchange if presents, then create a new one
func declareExchange(conn *amq.Connection, exchange, typee string) error {
	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	// err = ch.ExchangeDelete(exchange, false, false)
	// var amqpErr *amq.Error
	// if !(errors.As(err, &amqpErr) && amqpErr.Code == amq.NotFound) {
	// 	return err
	// }

	if err := ch.ExchangeDeclare(exchange, typee, false, false, false, false, nil); err != nil {
		return err
	}
	return err
}

// Only TCP packet sent confirmation
func publish(ctx context.Context, ch *amq.Channel, exchange, route string, msg []byte) {
	if err := ch.PublishWithContext(ctx, exchange, route, false, false, amq.Publishing{
		ContentType: "text/plain",
		Body:        msg,
	}); err != nil {
		panic(err)
	}
	fmt.Println("(TCP) MSG PUBLISHED TO: ", exchange, ", ROUTE:", route)
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

	// ---- // Confirmation Mode for publish

	// confirming mode
	if err := ch.Confirm(false); err != nil {
		panic(err)
	}

	// publish notifications are comming here
	confirms := ch.NotifyPublish(make(chan amq.Confirmation, 1))

	go func() {
		fmt.Println("I AM INSIDE THE CONFIRMATION GOROUTINE.")
		for c := range confirms {
			if c.Ack {
				fmt.Println(c.DeliveryTag, "published!!!")
			} else {
				fmt.Println(c.DeliveryTag, "not published ;__) !")
			}
		}
	}()

	// --- //

	// exchange, type
	ex := "logs"
	exType := "direct"

	if err := declareExchange(conn, ex, exType); err != nil {
		panic(err)
	}

	route1 := "logs.info"
	queue1 := "info-logs"
	if err := declareQueue(conn, ex, queue1, route1); err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// --- // Publish Channel must be same as Confirm
	for i := 0; i <= 20; i++ {
		publish(ctx, ch, ex, route1, []byte("i am"+strconv.Itoa(i)))
	}

	time.Sleep(30 * time.Second)
}
