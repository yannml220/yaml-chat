package user

import (
	"context"
	"github.com/yannml220/chat_agent_app/internal/pkg/types"
)

type UserService interface {
	GetUserById(ctx context.Context, id string) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	UpdateUser(ctx context.Context, user *User, id string) error
	DeleteUser(ctx context.Context, id string) error
}

type UserServiceImpl struct {
	types.Service
	repo UserRepo
}

func NewService(repo UserRepo) UserService {

	return &UserServiceImpl{
		types.Service{},
		repo,
	}

}

func (us UserServiceImpl) GetUserById(ctx context.Context, id string) (*User, error) {

	return us.repo.GetUserById(ctx, id)

}

func (us UserServiceImpl) GetUserByEmail(ctx context.Context, email string) (*User, error) {

	return us.repo.GetUserByEmail(ctx, email)

}

func (us UserServiceImpl) UpdateUser(ctx context.Context, user *User, id string) error {
	return us.repo.UpdateUser(ctx, user, id)

}

func (us UserServiceImpl) DeleteUser(ctx context.Context, id string) error {
	return us.repo.DeleteUser(ctx, id)

}
