package user

import (
	"time"
)

type User struct {
	Id                string     `json:"id"`
	Username          *string    `json:"username"`
	Name              *string    `json:"name"`
	ProfilePictureUrl *string    `json:"profile_picture_url"`
	Email             *string    `json:"email"`
	Confirmed         *bool      `json:"confirmed"`
	CreatedAt         *time.Time `json:"created_at,omitempty"`
	UpdatedAt         *time.Time `json:"updated_at,omitempty"`
	DeletedAt         *time.Time `json:"deleted_at,omitempty"`
}
