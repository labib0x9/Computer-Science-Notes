package user

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockUserRepo struct {
	users map[int]*User
	err   error
}

func (m *mockUserRepo) FindById(id int) (*User, error) {
	if m.err != nil {
		return nil, m.err
	}
	u, ok := m.users[id]
	if !ok {
		return nil, nil
	}
	return u, nil
}

func TestUserService_GetUserName(t *testing.T) {
	mockRepo := &mockUserRepo{
		users: map[int]*User{
			1 : {Id: 1, Name: "A"},
			2 : {Id: 2, Name: "B"},
			3 : {Id: 3, Name: "C"},
			4 : {Id: 4, Name: "D"},
		},
	}

	service := NewUserService(mockRepo)

	name, err := service.GetUserName(1)
	assert.NoError(t, err)
	assert.Equal(t, "A", name)

	name, err = service.GetUserName(6)
	assert.Error(t, err)
	assert.Equal(t, "", name)

	mockRepo.err = errors.New("Down")
	name, err = service.GetUserName(1)
	assert.Error(t, err)
	assert.Equal(t, "", name)
}
