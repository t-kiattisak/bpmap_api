package repository

import (
	"context"
	"pbmap_api/src/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id uuid.UUID) error
	FindAll(ctx context.Context) ([]domain.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	return GetDB(ctx, r.db).Create(user).Error
}

func (r *userRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	var user domain.User
	err := GetDB(ctx, r.db).Preload("SocialAccounts").
		Preload("SpecialCredential").
		Preload("Devices").
		Preload("Sessions").
		First(&user, "id = ?", id).Error
	return &user, err
}

func (r *userRepository) Update(ctx context.Context, user *domain.User) error {
	return GetDB(ctx, r.db).Model(user).Updates(user).Error
}

func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return GetDB(ctx, r.db).Delete(&domain.User{}, "id = ?", id).Error
}

func (r *userRepository) FindAll(ctx context.Context) ([]domain.User, error) {
	var users []domain.User
	err := GetDB(ctx, r.db).Find(&users).Error
	return users, err
}
