package main

import "log"

type ConsoleLogger struct {}

func (cl ConsoleLogger) Log(msg string) {
	log.Println(msg)
}

type Logger interface {
	Log(msg string)
}