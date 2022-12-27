package services

import (
	"github.com/ttanik/bookstore_users-api/domain/users"
	"github.com/ttanik/bookstore_users-api/utils/crypto_utils"
	"github.com/ttanik/bookstore_utils-go/rest_errors"
)

var (
	UsersService usersServiceInterface = &usersService{}
)

type usersService struct{}

type usersServiceInterface interface {
	CreateUser(users.User) (*users.User, *rest_errors.RestErr)
	GetUser(int64) (*users.User, *rest_errors.RestErr)
	DeleteUser(int64) *rest_errors.RestErr
	UpdateUser(bool, users.User) (*users.User, *rest_errors.RestErr)
	SearchUser(string) (users.Users, *rest_errors.RestErr)
	LoginUser(users.LoginRequest) (*users.User, *rest_errors.RestErr)
}

func (s *usersService) CreateUser(user users.User) (*users.User, *rest_errors.RestErr) {
	if err := user.Validate(); err != nil {
		return nil, err
	}
	user.Status = users.StatusActive
	user.Password = crypto_utils.GetMd5(user.Password)
	if err := user.Save(); err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *usersService) GetUser(userId int64) (*users.User, *rest_errors.RestErr) {
	result := users.User{Id: userId}
	if err := result.Get(); err != nil {
		return nil, err
	}
	return &result, nil
}
func (s *usersService) DeleteUser(userId int64) *rest_errors.RestErr {
	deleteUser := users.User{Id: userId}
	if err := deleteUser.Delete(); err != nil {
		return err
	}
	return nil
}
func (s *usersService) UpdateUser(isPartial bool, user users.User) (*users.User, *rest_errors.RestErr) {
	err := user.Get()
	if err != nil {
		return nil, err
	}
	current := &user
	if err := current.Validate(); err != nil {
		return nil, err
	}
	if isPartial {
		if user.FirstName != "" {
			current.FirstName = user.FirstName
		}
		if user.LastName != "" {
			current.LastName = user.LastName
		}
		if user.Email != "" {
			current.Email = user.Email
		}
	} else {
		current.FirstName = user.FirstName
		current.LastName = user.LastName
		current.Email = user.Email
	}
	if err := current.Update(); err != nil {
		return nil, err
	}
	return current, nil
}

func (s *usersService) SearchUser(status string) (users.Users, *rest_errors.RestErr) {
	dao := &users.User{}
	return dao.FindByStatus(status)
}

func (s *usersService) LoginUser(loginRequest users.LoginRequest) (*users.User, *rest_errors.RestErr) {

	user := &users.User{
		Email:    loginRequest.Email,
		Password: crypto_utils.GetMd5(loginRequest.Password),
	}
	err := user.FindByEmailAndPassword()
	if err != nil {
		return nil, err
	}
	return user, nil
}
