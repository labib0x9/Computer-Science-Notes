package main

import (
	"fmt"
	"sync"
	"time"
)

var adminSender, userSender chan string
var wg sync.WaitGroup

func AdminSend(msg string) {
	defer wg.Done()
	adminSender <- msg
}

func UserSender(msg string) {
	defer wg.Done()
	userSender <- msg
}

func Check() {
	defer wg.Done()
	counter := 0
free:
	for {
		select {
		case msg := <-adminSender:
			counter++
			fmt.Println("Admin: ", msg)
		case msg, ok := <-userSender:
			if !ok {
				// channel closed
				continue
			}
			counter++
			fmt.Println("User: ", msg)
		case <-time.After(5 * time.Millisecond): // If no channel is receiving data in 5 milisecond, then it would return
			fmt.Println("Timeout")
			break free
			// default:
			// 	fmt.Println("DEFAULT: ", counter)
			// 	break free
		}
	}
}

func StayIdle() {
	defer wg.Done()
	timer := time.NewTimer(5 * time.Millisecond) // Receiving time is 5 mili sec
	defer timer.Stop()

	for {
		select {
		case msg := <-adminSender:
			fmt.Println("Admin: ", msg)
			timer.Reset(3 * time.Millisecond) // Receiving time is set to 3 mili sec
		case msg := <-userSender:
			fmt.Println("User: ", msg)
			timer.Reset(2 * time.Millisecond) // Receiving time is set to 2 mili sec
		case <-timer.C: // After receiving time, it will return
			fmt.Println("Timeout")
			return
		}
	}
}

func main() {

	adminSender = make(chan string)
	userSender = make(chan string)

	wg.Add(3)

	go AdminSend("Hello Client")
	go UserSender("Hello Server")
	// go Check()
	go StayIdle()

	wg.Wait()
}
