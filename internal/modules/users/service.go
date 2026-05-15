package users

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/example/gin-go-monolithic-boilerplate/internal/common/apperror"
	"github.com/example/gin-go-monolithic-boilerplate/internal/common/id"
)

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) List(ctx context.Context) ([]UserResponse, error) {
	users, err := s.repository.List(ctx)
	if err != nil {
		return nil, err
	}
	return toUserResponses(users), nil
}

func (s *Service) GetByID(ctx context.Context, userID string) (UserResponse, error) {
	user, err := s.repository.GetByID(ctx, strings.TrimSpace(userID))
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return UserResponse{}, apperror.NotFound("User not found.")
		}
		return UserResponse{}, err
	}
	return toUserResponse(user), nil
}

func (s *Service) Create(ctx context.Context, request CreateUserRequest) (UserResponse, error) {
	now := time.Now().UTC()
	user := User{ID: id.New(), Name: strings.TrimSpace(request.Name), Email: normalizeEmail(request.Email), CreatedAt: now, UpdatedAt: now}

	created, err := s.repository.Create(ctx, user)
	if err != nil {
		if errors.Is(err, ErrEmailTaken) {
			return UserResponse{}, apperror.Conflict("A user with this email already exists.")
		}
		return UserResponse{}, err
	}
	return toUserResponse(created), nil
}

func (s *Service) Update(ctx context.Context, userID string, request UpdateUserRequest) (UserResponse, error) {
	user, err := s.repository.GetByID(ctx, strings.TrimSpace(userID))
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return UserResponse{}, apperror.NotFound("User not found.")
		}
		return UserResponse{}, err
	}

	if request.Name != nil {
		user.Name = strings.TrimSpace(*request.Name)
	}
	if request.Email != nil {
		user.Email = normalizeEmail(*request.Email)
	}
	user.UpdatedAt = time.Now().UTC()

	updated, err := s.repository.Update(ctx, user)
	if err != nil {
		if errors.Is(err, ErrEmailTaken) {
			return UserResponse{}, apperror.Conflict("A user with this email already exists.")
		}
		if errors.Is(err, ErrUserNotFound) {
			return UserResponse{}, apperror.NotFound("User not found.")
		}
		return UserResponse{}, err
	}
	return toUserResponse(updated), nil
}

func (s *Service) Delete(ctx context.Context, userID string) error {
	err := s.repository.Delete(ctx, strings.TrimSpace(userID))
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return apperror.NotFound("User not found.")
		}
		return err
	}
	return nil
}
