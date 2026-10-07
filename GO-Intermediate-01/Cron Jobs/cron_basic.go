package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-co-op/gocron/v2"
)

/*
	Task Scheduler

	go get github.com/go-co-op/gocron/v2
*/

func main() {

	s, err := gocron.NewScheduler()
	if err != nil {
		panic(err)
	}

	// runs every 22 seconds
	jd := gocron.DurationJob(22 * time.Second)
	task := gocron.NewTask(func(a string) {
		log.Println(a)
	}, "Every 22 sec")

	// gocron.CronJob("*/10 * * * *", false)	// cron - every 10 minutes
	// gocron.DailyJob(1, gocron.NewAtTimes(gocron.NewAtTime(23, 10, 0), gocron.NewAtTime(5, 30,0)))	// daily, 1 = everyday, newattime = trigger
	// gocron.OneTimeJob(gocron.OneTimeJobStartDateTime(time.Now().Add(10 * time.Minute)))

	name := gocron.WithName("First CRON")

	j, err := s.NewJob(jd, task, name)
	if err != nil {
		panic(err)
	}

	log.Println(j.ID())

	s.Start()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	<-ctx.Done()
	s.Shutdown()
}
