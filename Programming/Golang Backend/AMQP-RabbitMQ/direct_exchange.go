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
	if err := ch.ExchangeDeclare(exchange, typee, false, false, false, false, nil); err != nil {
		return err
	}

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

func consume(conn *amq.Connection, queue string, consumer string) {
	ch, err := conn.Channel()
	if err != nil {
		fmt.Println(consumer, "ERROR:", err)
		return
	}
	defer ch.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	msgs, err := ch.ConsumeWithContext(ctx, queue, consumer, true, false, false, false, nil)

	for {
		select {
		case <-ctx.Done():
			fmt.Println(consumer, "RETURNING..")
			return
		case m := <-msgs:
			fmt.Println(consumer, "MSG:", string(m.Body))
		}
	}
}

func publish(ctx context.Context, conn *amq.Connection, exchange, route string, msg []byte) {
	ch, err := conn.Channel()
	if err != nil {
		fmt.Println(exchange, "ERROR:", err)
		return
	}
	defer ch.Close()
	if err := ch.PublishWithContext(ctx, exchange, route, false, false, amq.Publishing{
		ContentType: "text/plain",
		Body:        msg,
	}); err != nil {
		panic(err)
	}
	fmt.Println("MSG PUBLISHED TO: ", exchange, ", ROUTE:", route)
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

	// exchange, type
	ex := "logs"
	exType := "direct"

	route1 := "logs.info"
	queue1 := "info-logs"
	if err := declareQueue(ch, ex, exType, queue1, route1); err != nil {
		panic(err)
	}

	route2 := "logs.warning"
	queue2 := "warning-logs"
	if err := declareQueue(ch, ex, exType, queue2, route2); err != nil {
		panic(err)
	}

	go consume(conn, queue1, "INFO_CONSUMER")
	go consume(conn, queue2, "WARNING_CONSUMER")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	route3 := "logs.#"
	// one exchange -> routes to multiple queue (based on routing key)
	publish(ctx, conn, ex, route1, []byte("only be in queue= info logs"))
	publish(ctx, conn, ex, route2, []byte("only be in queue= warning logs"))
	publish(ctx, conn, ex, route3, []byte("only be in queue= none (lost)"))

	time.Sleep(30 * time.Second)
}
