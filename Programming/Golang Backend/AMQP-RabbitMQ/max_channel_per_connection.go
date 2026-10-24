package main

import (
	"fmt"
	"sync"
	"time"

	amq "github.com/rabbitmq/amqp091-go"
)

func main() {

	conn, err := amq.DialConfig(
		"amqp://guest:guest@localhost:5672",
		amq.Config{
			ChannelMax: 10,
			Recovery:   &amq.Recovery{},
		})
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	var wg sync.WaitGroup
	for i := range 20 {
		wg.Add(1)
		go func(id int, wg *sync.WaitGroup) {
			defer wg.Done()
			ch, err := conn.Channel()
			if err != nil {
				fmt.Println("ID:", id, "Error", err, "Coonection", conn.IsClosed())
				return
			}
			fmt.Println("ID:", id, "Coonection", conn.IsClosed())
			defer ch.Close()
			time.Sleep(10 * time.Second)
		}(i, &wg)
	}

	wg.Wait()
}
