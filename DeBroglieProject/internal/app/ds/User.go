package ds

import (
	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Email       string    `gorm:"type:varchar(25);unique;not null" json:"email"`
	Name        string    `gorm:"type:varchar(50);not null" json:"name"`
	Password    string    `gorm:"type:varchar(255);not null" json:"-"`
	IsProfessor bool      `gorm:"type:boolean;default:false" json:"is_professor"`
}
