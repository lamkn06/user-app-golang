package service

import (
	"github.com/google/uuid"

	"github.com/lamkn06/user-app-golang.git/internal/repository"

	"github.com/lamkn06/user-app-golang.git/pkg/api/request"
	"github.com/lamkn06/user-app-golang.git/pkg/api/response"
)

type UserService interface {
	GetUsers(listReq request.ListRequest) (response.ListResponse[response.NewUserResponse], error)
	NewUser(user request.NewUserRequest) (response.NewUserResponse, error)
	GetUserById(id uuid.UUID) (response.NewUserResponse, error)
}

type DefaultUserService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return &DefaultUserService{userRepository: userRepository}
}

func (s *DefaultUserService) GetUsers(listReq request.ListRequest) (response.ListResponse[response.NewUserResponse], error) {
	// Get total count
	total, err := s.userRepository.GetUsersCount()
	if err != nil {
		return response.ListResponse[response.NewUserResponse]{}, err
	}

	// Get paginated users
	users, err := s.userRepository.GetUsers(listReq.GetOffset(), listReq.GetLimit())
	if err != nil {
		return response.ListResponse[response.NewUserResponse]{}, err
	}

	// Convert to response format
	var responses []response.NewUserResponse
	for _, user := range users {
		responses = append(responses, response.NewUserResponse{
			ID:    user.Id,
			Name:  user.Name,
			Email: user.Email,
		})
	}

	// Create list response with metadata
	return response.NewListResponse(responses, total, listReq.GetPage(), listReq.GetLimit()), nil
}

func (s *DefaultUserService) NewUser(user request.NewUserRequest) (response.NewUserResponse, error) {
	entity := repository.UserEntity{
		Id:    uuid.New(),
		Name:  user.Name,
		Email: user.Email,
	}

	newUser, err := s.userRepository.InsertUser(entity)
	if err != nil {
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
