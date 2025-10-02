package repository

import (
	"DeBroglieProject/internal/app/ds"

	"github.com/google/uuid"
)

func (r *Repository) GetUserByUUID(uuid uuid.UUID) (ds.User, error) {
	var u ds.User
	err := r.db.Where("id = ?", uuid).First(&u).Error
	return u, err
}

func (r *Repository) GetUserByEmail(email string) (ds.User, error) {
	var u ds.User
	err := r.db.Where("email = ?", email).First(&u).Error
	return u, err
}

func (r *Repository) UpdateUserByUUID(uuid uuid.UUID, user ds.User) error {
	updates := map[string]interface{}{"name": user.Name}
	updates["is_moderator"] = user.IsModerator
	return r.db.Model(&ds.User{}).Where("id = ?", uuid).Updates(updates).Error
}

func (r *Repository) Register(user *ds.User) error {
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}
	return r.db.Create(user).Error
}
