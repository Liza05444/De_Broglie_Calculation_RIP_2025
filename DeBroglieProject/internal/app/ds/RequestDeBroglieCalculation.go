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
	ID          uint          `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string        `gorm:"type:varchar(100);not null" json:"name"`
	Status      RequestStatus `gorm:"type:varchar(20);not null;check:status IN ('черновик','удалён','сформирован','завершён','отклонён')" json:"status"`
	CreatedAt   time.Time     `gorm:"not null" json:"created_at"`
	FormedAt    sql.NullTime  `gorm:"default:null" json:"formed_at,omitempty"`
	CompletedAt sql.NullTime  `gorm:"default:null" json:"completed_at,omitempty"`
	CreatorID   uint          `gorm:"not null" json:"creator_id"`
	ModeratorID sql.NullInt64 `json:"moderator_id,omitempty"`

	Creator   User `gorm:"foreignKey:CreatorID" json:"creator,omitempty"`
	Moderator User `gorm:"foreignKey:ModeratorID" json:"moderator,omitempty"`
}
