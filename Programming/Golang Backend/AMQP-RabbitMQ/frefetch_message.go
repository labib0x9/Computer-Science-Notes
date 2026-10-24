func consume(conn *amq.Connection, queue string, consumer string) {
	ch, err := conn.Channel()
	if err != nil {
		fmt.Println(consumer, "ERROR:", err)
		return
	}
	defer ch.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	// how many concurrent msg to handle at once
	// need to undestand the global parameter
	err = ch.Qos(10, 0, false) // per consumer limit
	// ch.Qos(10, 0, true)  // per channel limit
	if err != nil {
		fmt.Println(consumer, "ERROR:", err)
		return
	}

	msgs, err := ch.ConsumeWithContext(ctx, queue, consumer, true, false, false, false, nil)
	if err != nil {
		fmt.Println(consumer, "ERROR:", err)
		return
	}

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
