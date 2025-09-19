package ds

import (
	"database/sql"
)

type Particle struct {
	ID          uint           `gorm:"primaryKey"`
	Name        string         `gorm:"type:varchar(50);not null"`
	Mass        float64        `gorm:"type:numeric;not null"`
	Image       sql.NullString `gorm:"type:varchar(150);default:null"`
	Description sql.NullString `gorm:"type:text;default:null"`
	IsDeleted   bool           `gorm:"type:boolean;not null;default:false"`
}
