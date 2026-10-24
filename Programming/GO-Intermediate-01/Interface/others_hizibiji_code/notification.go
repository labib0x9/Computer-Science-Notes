package main

import "fmt"

/**

**/

type User struct {
	Name string
}

type Admin struct {
	Name string
	Id int
}

type Email struct {
	EmailAddress string
}

type Notification interface {
	Notify(string)
}

func Notifier(n Notification, msg string) {

	// Type assertion - 01
	if admin, ok := n.(Admin); ok {
		msg += fmt.Sprintf(" : Admin id %d", admin.Id)
	}

	n.Notify(msg)

	// Type switch - 01
	switch v := n.(type) {
	case Admin:
		fmt.Println("Admin ::", v.Name)
	case User:
		fmt.Println("User ::", v.Name)
	case Email:
		fmt.Println("Email Address ::", v.EmailAddress)
	}
}

func (u User) Notify(msg string) {
	fmt.Printf("[ User ] %s : %s\n", u.Name, msg)
}

func (a Admin) Notify(msg string) {
	fmt.Printf("[ Admin ] %s : %s\n", a.Name, msg)
}

func (e Email) Notify(msg string) {
	fmt.Printf("%s is send to %s\n", msg, e.EmailAddress)
}

func NotificationManagement() {
	user1 := User{Name: "Labib"}
	admin1 := Admin{Name: "Admin-01", Id: 1}
	email1 := Email{EmailAddress: "labib@gmail.com"}

	Notifier(user1, "I am trying to login")
	Notifier(admin1, "Server is down !! Check logs")
	Notifier(email1, "Are you trying to login, Check notification")
}