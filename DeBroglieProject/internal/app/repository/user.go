package repository

import (
	"DeBroglieProject/internal/app/ds"
	"errors"
)

func (r *Repository) CreateUser(user ds.User) (ds.User, error) {
	if user.Name == "" {
		user.Name = "Пользователь"
	}
	err := r.db.Create(&user).Error
	return user, err
}

func (r *Repository) GetUserByID(id uint) (ds.User, error) {
	var u ds.User
	err := r.db.First(&u, id).Error
	return u, err
}

func (r *Repository) GetUserByEmail(email string) (ds.User, error) {
	var u ds.User
	err := r.db.Where("email = ?", email).First(&u).Error
	return u, err
}

func (r *Repository) UpdateUser(id uint, user ds.User) error {
	updates := map[string]interface{}{"name": user.Name}
	updates["is_moderator"] = user.IsModerator
	return r.db.Model(&ds.User{}).Where("id = ?", id).Updates(updates).Error
}

func (r *Repository) CheckCredentials(email, password string) (ds.User, error) {
	u, err := r.GetUserByEmail(email)
	if err != nil {
		return ds.User{}, err
	}
	if u.Password != password {
		return ds.User{}, errors.New("invalid credentials")
	}
	return u, nil
}
