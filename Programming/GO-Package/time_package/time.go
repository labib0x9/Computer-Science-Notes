package main

import "time"

func main() {
	now := time.Now()

	now.Year()
	now.Month()
	now.Day()
	now.Date()
	now.Hour()
	now.Minute()
	now.Second()
	now.Weekday()

	now.Add(2 * time.Hour)
	now.AddDate(1, 0, 0)

	time.Sleep(1 * time.Microsecond)

	time.NewTimer(5 * time.Microsecond)
	time.NewTicker(5 * time.Microsecond)

	now.Format("2006-01-02 15:04:05")
	// time.Parse("2006-01-02 15:04:05", )

	time.Since(now)
}