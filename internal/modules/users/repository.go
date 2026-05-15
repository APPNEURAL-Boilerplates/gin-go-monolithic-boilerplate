package users

import (
	"context"
	"errors"
	"sort"
	"strings"
	"sync"
)

var ErrUserNotFound = errors.New("user not found")
var ErrEmailTaken = errors.New("email is already taken")

type Repository interface {
	List(ctx context.Context) ([]User, error)
	GetByID(ctx context.Context, id string) (User, error)
	GetByEmail(ctx context.Context, email string) (User, error)
	Create(ctx context.Context, user User) (User, error)
	Update(ctx context.Context, user User) (User, error)
	Delete(ctx context.Context, id string) error
}

type InMemoryRepository struct {
	mu         sync.RWMutex
	users      map[string]User
	emailIndex map[string]string
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{users: map[string]User{}, emailIndex: map[string]string{}}
}

func (r *InMemoryRepository) List(ctx context.Context) ([]User, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	items := make([]User, 0, len(r.users))
	for _, user := range r.users {
		items = append(items, user)
	}
	sort.Slice(items, func(i, j int) bool { return items[i].CreatedAt.Before(items[j].CreatedAt) })
	return items, nil
}

func (r *InMemoryRepository) GetByID(ctx context.Context, id string) (User, error) {
	select {
	case <-ctx.Done():
		return User{}, ctx.Err()
	default:
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	user, ok := r.users[id]
	if !ok {
		return User{}, ErrUserNotFound
	}
	return user, nil
}

func (r *InMemoryRepository) GetByEmail(ctx context.Context, email string) (User, error) {
	select {
	case <-ctx.Done():
		return User{}, ctx.Err()
	default:
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	id, ok := r.emailIndex[normalizeEmail(email)]
	if !ok {
		return User{}, ErrUserNotFound
	}
	user, ok := r.users[id]
	if !ok {
		return User{}, ErrUserNotFound
	}
	return user, nil
}

func (r *InMemoryRepository) Create(ctx context.Context, user User) (User, error) {
	select {
	case <-ctx.Done():
		return User{}, ctx.Err()
	default:
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	email := normalizeEmail(user.Email)
	if _, exists := r.emailIndex[email]; exists {
		return User{}, ErrEmailTaken
	}

	r.users[user.ID] = user
	r.emailIndex[email] = user.ID
	return user, nil
}

func (r *InMemoryRepository) Update(ctx context.Context, user User) (User, error) {
	select {
	case <-ctx.Done():
		return User{}, ctx.Err()
	default:
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	current, ok := r.users[user.ID]
	if !ok {
		return User{}, ErrUserNotFound
	}

	oldEmail := normalizeEmail(current.Email)
	newEmail := normalizeEmail(user.Email)
	if oldEmail != newEmail {
		if existingID, exists := r.emailIndex[newEmail]; exists && existingID != user.ID {
			return User{}, ErrEmailTaken
		}
		delete(r.emailIndex, oldEmail)
		r.emailIndex[newEmail] = user.ID
	}

	r.users[user.ID] = user
	return user, nil
}

func (r *InMemoryRepository) Delete(ctx context.Context, id string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	user, ok := r.users[id]
	if !ok {
		return ErrUserNotFound
	}

	delete(r.users, id)
	delete(r.emailIndex, normalizeEmail(user.Email))
	return nil
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}
