package service

import (
	"errors"

	"github.com/google/uuid"
	"github.com/lamkn06/user-app-golang.git/internal/repository"
	"github.com/lamkn06/user-app-golang.git/pkg/api/request"
	"github.com/lamkn06/user-app-golang.git/pkg/api/response"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	SignUp(req request.SignUpRequest) (response.AuthResponse, error)
	SignIn(req request.SignInRequest) (response.AuthResponse, error)
	SignOut(token string) (response.SignOutResponse, error)
}

type DefaultAuthService struct {
	userRepository repository.UserRepository
	jwtService     JWTService
}

func NewAuthService(userRepository repository.UserRepository, jwtService JWTService) AuthService {
	return &DefaultAuthService{
		userRepository: userRepository,
		jwtService:     jwtService,
	}
}

func (s *DefaultAuthService) SignUp(req request.SignUpRequest) (response.AuthResponse, error) {
	// Check if user already exists
	_, err := s.userRepository.GetUserByEmail(req.Email)
	if err == nil {
		return response.AuthResponse{}, errors.New("user already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return response.AuthResponse{}, err
	}

	// Create user
	user := repository.UserEntity{
		Id:       uuid.New(),
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	createdUser, err := s.userRepository.InsertUser(user)
	if err != nil {
		return response.AuthResponse{}, err
	}

	// Generate tokens
	token, err := s.jwtService.GenerateToken(createdUser.Id, createdUser.Email)
	if err != nil {
		return response.AuthResponse{}, err
	}

	refreshToken, err := s.jwtService.GenerateRefreshToken(createdUser.Id)
	if err != nil {
		return response.AuthResponse{}, err
	}

	return response.AuthResponse{
		Token:        token,
		RefreshToken: refreshToken,
		User: response.UserResponse{
			ID:    createdUser.Id.String(),
			Name:  createdUser.Name,
			Email: createdUser.Email,
		},
	}, nil
}

func (s *DefaultAuthService) SignIn(req request.SignInRequest) (response.AuthResponse, error) {
	// Get user by email
	user, err := s.userRepository.GetUserByEmail(req.Email)
	if err != nil {
		return response.AuthResponse{}, errors.New("invalid credentials")
	}

	// Check password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return response.AuthResponse{}, errors.New("invalid credentials")
	}

	// Generate tokens
	token, err := s.jwtService.GenerateToken(user.Id, user.Email)
	if err != nil {
		return response.AuthResponse{}, err
	}

	refreshToken, err := s.jwtService.GenerateRefreshToken(user.Id)
	if err != nil {
		return response.AuthResponse{}, err
	}

	return response.AuthResponse{
		Token:        token,
		RefreshToken: refreshToken,
		User: response.UserResponse{
			ID:    user.Id.String(),
			Name:  user.Name,
			Email: user.Email,
		},
	}, nil
}

func (s *DefaultAuthService) SignOut(token string) (response.SignOutResponse, error) {
	// In a real application, you would blacklist the token
	// For now, we just return success
	return response.SignOutResponse{
		Message: "Successfully signed out",
	}, nil
}
