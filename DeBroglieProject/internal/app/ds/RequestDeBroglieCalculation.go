package ds

import (
	"database/sql"
	"time"
)

type RequestStatus string

const (
	RequestStatusDraft     RequestStatus = "черновик"
	RequestStatusDeleted   RequestStatus = "удалён"
	RequestStatusFormed    RequestStatus = "сформирован"
	RequestStatusCompleted RequestStatus = "завершён"
	RequestStatusRejected  RequestStatus = "отклонён"
)

type RequestDeBroglieCalculation struct {
	ID          uint          `gorm:"primaryKey;autoIncrement"`
	Name        string        `gorm:"type:varchar(100);not null"`
	Status      RequestStatus `gorm:"type:varchar(20);not null;check:status IN ('черновик','удалён','сформирован','завершён','отклонён')"`
	CreatedAt   time.Time     `gorm:"not null"`
	FormedAt    sql.NullTime  `gorm:"default:null"`
	CompletedAt sql.NullTime  `gorm:"default:null"`
	CreatorID   uint          `gorm:"not null"`
	ModeratorID sql.NullInt64

	Creator   User `gorm:"foreignKey:CreatorID"`
	Moderator User `gorm:"foreignKey:ModeratorID"`
}
