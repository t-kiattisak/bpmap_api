package usecase

import (
	"context"
	"pbmap_api/src/domain"
	"pbmap_api/src/internal/repository"

	"github.com/google/uuid"
)

type UserUsecase interface {
	CreateUser(ctx context.Context, user *domain.User) error
	GetUser(ctx context.Context, id uuid.UUID) (*domain.User, error)
	UpdateUser(ctx context.Context, user *domain.User) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
	ListUsers(ctx context.Context) ([]domain.User, error)
}

type userUsecase struct {
	userRepo repository.UserRepository
}

func NewUserUsecase(userRepo repository.UserRepository) UserUsecase {
	return &userUsecase{userRepo}
}

func (u *userUsecase) CreateUser(ctx context.Context, user *domain.User) error {
	return u.userRepo.Create(ctx, user)
}

func (u *userUsecase) GetUser(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	return u.userRepo.FindByID(ctx, id)
}

func (u *userUsecase) UpdateUser(ctx context.Context, user *domain.User) error {
	return u.userRepo.Update(ctx, user)
}

func (u *userUsecase) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return u.userRepo.Delete(ctx, id)
}

func (u *userUsecase) ListUsers(ctx context.Context) ([]domain.User, error) {
	return u.userRepo.FindAll(ctx)
}
