package users

import "time"

type CreateUserRequest struct {
	Name  string `json:"name" binding:"required,min=2,max=100"`
	Email string `json:"email" binding:"required,email"`
}

type UpdateUserRequest struct {
	Name  *string `json:"name" binding:"omitempty,min=2,max=100"`
	Email *string `json:"email" binding:"omitempty,email"`
}

type UserResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func toUserResponse(user User) UserResponse {
	return UserResponse{ID: user.ID, Name: user.Name, Email: user.Email, CreatedAt: user.CreatedAt, UpdatedAt: user.UpdatedAt}
}

func toUserResponses(users []User) []UserResponse {
	result := make([]UserResponse, 0, len(users))
	for _, user := range users {
		result = append(result, toUserResponse(user))
	}
	return result
}
