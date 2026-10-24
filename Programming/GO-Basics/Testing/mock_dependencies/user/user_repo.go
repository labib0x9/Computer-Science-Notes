package user

type User struct {
	Id   int
	Name string
}

// Why interface is needed ??
type UserRepository interface {
	FindById(id int) (*User, error)
}