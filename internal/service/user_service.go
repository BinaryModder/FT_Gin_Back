package service

import (
	"errors"
	"github.com/BinaryModder/FT_Gin_Back.git/internal/model"
	"github.com/BinaryModder/FT_Gin_Back.git/internal/repository"
	"strings"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) RegisterUser(name, email string, age int) (*model.User, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, errors.New("name is required")
	}

	email = strings.TrimSpace(email)
	if email == "" {
		return nil, errors.New("email is required")
	}
	existingUser, _ := s.repo.GetUserByEmail(email)
	if existingUser != nil {
		return nil, errors.New("email already registered")
	}
	user := &model.User{
		Name:  name,
		Email: email,
		Age:   age,
	}
	if err := s.repo.CreateUser(user); err != nil {
		return nil, errors.New("failed to create user")
	}

	return user, nil
}
func (s *UserService) GetUser(id uint) (*model.User, error) {
	if id == 0 {
		return nil, errors.New("invalid user id")
	}

	user, err := s.repo.GetUserByID(id)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}

func (s *UserService) GetAllUsers() ([]model.User, error) {
	users, err := s.repo.GetAllUsers()
	if err != nil {
		return nil, errors.New("failed to fetch users")
	}
	return users, nil
}

func (s *UserService) UpdateUserInfo(id uint, name, email string, age int) (*model.User, error) {
	if id == 0 {
		return nil, errors.New("invalid user id")
	}

	user, err := s.repo.GetUserByID(id)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if name != "" {
		user.Name = strings.TrimSpace(name)
	}

	if email != "" {
		email = strings.TrimSpace(email)
		if email != user.Email {
			existingUser, _ := s.repo.GetUserByEmail(email)
			if existingUser != nil {
				return nil, errors.New("email already in use")
			}
			user.Email = email
		}
	}
	user.Age = age

	if err := s.repo.UpdateUser(user); err != nil {
		return nil, errors.New("failed to update user")
	}

	return user, nil
}

func (s *UserService) DeleteUser(id uint) error {
	if id == 0 {
		return errors.New("invalid user id")
	}

	_, err := s.repo.GetUserByID(id)
	if err != nil {
		return errors.New("user not found")
	}

	return s.repo.DeleteUser(id)
}
