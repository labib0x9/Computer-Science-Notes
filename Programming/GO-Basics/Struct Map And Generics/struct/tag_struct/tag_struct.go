package main

// What does omitempty, string, - means
// And how do they work ??

type User struct {
	Username string `json:"username"`
	Email    string `json:"email,omitempty"`
	Id       int    `json:"id,string"`
	Password string `json:"-"` // Exclude password field from json output
}

func main() {

}