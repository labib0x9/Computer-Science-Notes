package main

func main() {
	logger := ConsoleLogger{}
	service := NewUserService(logger)
	service.CreateUser("Labib")
}