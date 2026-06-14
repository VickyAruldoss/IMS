package model

import "time"

type Member struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateMemberRequest struct {
	Name  string `json:"name"  binding:"required"`
	Email string `json:"email" binding:"required,email"`
	Role  string `json:"role"  binding:"required"`
}

type UpdateMemberRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"  binding:"omitempty,email"`
	Role  string `json:"role"`
}
