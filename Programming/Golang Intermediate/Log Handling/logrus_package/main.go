package main

import (
	"os"

	"github.com/sirupsen/logrus"
)

func main() {

	logger := logrus.New()

	logger.SetLevel(logrus.InfoLevel)

	logger.SetFormatter(&logrus.JSONFormatter{})

	file, err := os.OpenFile("main.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		logger.Info("failed to open log file")
	} else {
		logger.SetOutput(file)
	}

	logger.Debug("Debug message")
	logger.Info("Info message")
	logger.Warn("Warning message")
	logger.Error("Error message")
	// logger.Fatal("Fatal message")
	// logger.Panic("Panic message")

	logger.WithFields(logrus.Fields{
		"user": "labib",
		"id":   42,
	}).Info("User logged in")

}
