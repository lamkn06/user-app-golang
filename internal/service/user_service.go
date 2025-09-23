package service

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/lamkn06/user-app-golang.git/internal/repository"

	"github.com/lamkn06/user-app-golang.git/pkg/api/request"
	"github.com/lamkn06/user-app-golang.git/pkg/api/response"
)

type UserService interface {
	GetUsers() ([]response.NewUserResponse, error)
	NewUser(user request.NewUserRequest) (response.NewUserResponse, error)
	GetUserById(id uuid.UUID) (response.NewUserResponse, error)
}

type DefaultUserService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return &DefaultUserService{userRepository: userRepository}
}

func (s *DefaultUserService) GetUsers() ([]response.NewUserResponse, error) {
	users, err := s.userRepository.GetUsers()
	if err != nil {
		return []response.NewUserResponse{}, err
	}

	if len(users) == 0 {
		return []response.NewUserResponse{}, nil
	}

	var responses []response.NewUserResponse
	for _, user := range users {
		responses = append(responses, response.NewUserResponse{
			ID:    user.Id,
			Name:  user.Name,
			Email: user.Email,
		})
	}
	return responses, nil
}

func (s *DefaultUserService) NewUser(user request.NewUserRequest) (response.NewUserResponse, error) {
	entity := repository.UserEntity{
		Id:    uuid.New(),
		Name:  user.Name,
		Email: user.Email,
	}

	newUser, err := s.userRepository.InsertUser(entity)
	if err != nil {
		fmt.Println("Error inserting user:", err)
		return response.NewUserResponse{}, err
	}

	return response.NewUserResponse{
		ID:    newUser.Id,
		Name:  newUser.Name,
		Email: newUser.Email,
	}, nil

}

func (s *DefaultUserService) GetUserById(id uuid.UUID) (response.NewUserResponse, error) {
	user, err := s.userRepository.GetUserById(id)
	if err != nil {
		return response.NewUserResponse{}, err
	}
	return response.NewUserResponse{
		ID:    user.Id,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}
