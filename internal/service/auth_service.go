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
	SignUp(req request.SignUpRequest) (response.SignUpResponse, error)
	SignIn(req request.SignInRequest) (response.SignInResponse, error)
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

func (s *DefaultAuthService) SignUp(req request.SignUpRequest) (response.SignUpResponse, error) {
	// Check if user already exists
	_, err := s.userRepository.GetUserByEmail(req.Email)
	if err == nil {
		return response.SignUpResponse{}, errors.New("user already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return response.SignUpResponse{}, err
	}

	// Create user
	user := repository.UserEntity{
		Id:       uuid.New(),
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	createdUser, err := s.userRepository.InsertUser(user)
	if err != nil {
		return response.SignUpResponse{}, err
	}

	return response.SignUpResponse{
		ID:    createdUser.Id.String(),
		Email: createdUser.Email,
	}, nil
}

func (s *DefaultAuthService) SignIn(req request.SignInRequest) (response.SignInResponse, error) {
	// Get user by email
	user, err := s.userRepository.GetUserByEmail(req.Email)
	if err != nil {
		return response.SignInResponse{}, errors.New("invalid credentials")
	}

	// Check password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return response.SignInResponse{}, errors.New("invalid credentials")
	}

	// Generate tokens
	token, err := s.jwtService.GenerateToken(user.Id, user.Email)
	if err != nil {
		return response.SignInResponse{}, err
	}

	refreshToken, err := s.jwtService.GenerateRefreshToken(user.Id)
	if err != nil {
		return response.SignInResponse{}, err
	}

	return response.SignInResponse{
		Token:        token,
		RefreshToken: refreshToken,
		User: response.NewUserResponse{
			ID:    user.Id.String(),
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
