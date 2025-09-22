package ds

import (
	"database/sql"
)

type Particle struct {
	ID          uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string         `gorm:"type:varchar(50);not null" json:"name"`
	Mass        float64        `gorm:"type:numeric;not null" json:"mass"`
	Image       sql.NullString `gorm:"type:varchar(50);default:null" json:"image,omitempty"`
	Description sql.NullString `gorm:"type:text;default:null" json:"description,omitempty"`
	IsDeleted   bool           `gorm:"type:boolean;not null;default:false" json:"-"`
}
