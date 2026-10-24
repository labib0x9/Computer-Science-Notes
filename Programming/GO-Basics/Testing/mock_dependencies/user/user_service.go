package user

import "errors"

// Why using interfaces inside a struct
type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

// To test this method
func (u *UserService) GetUserName(id int) (string, error) {
	user, err := u.repo.FindById(id)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errors.New("user not found")
	}
	return user.Name, nil
}