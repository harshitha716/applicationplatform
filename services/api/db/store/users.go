package store

import (
	"context"

	"github.com/Zampfi/application-platform/services/api/db/models"
)

type UserStore interface {
	GetUserById(ctx context.Context, userId string) (*models.User, error)
	GetUsersAll(ctx context.Context) ([]models.User, error)
}

func (s *appStore) GetUserById(ctx context.Context, userId string) (*models.User, error) {

	var user models.User
	err := s.client.WithContext(ctx).Model(user).Where("user_id = ?", userId).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *appStore) GetUsersAll(ctx context.Context) ([]models.User, error) {
	var users []models.User
	err := s.client.WithContext(ctx).Model(users).Order("email ASC").Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}
