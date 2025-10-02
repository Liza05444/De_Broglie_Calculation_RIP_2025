package ds

import (
	"time"
	"github.com/google/uuid"
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
	ID          uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string     `gorm:"type:varchar(100);not null" json:"name"`
	Status      RequestStatus `gorm:"type:varchar(20);not null;check:status IN ('черновик','удалён','сформирован','завершён','отклонён')" json:"status"`
	CreatedAt   time.Time  `gorm:"not null" json:"created_at"`
	FormedAt    *time.Time `gorm:"default:null" json:"formed_at,omitempty"`
	CompletedAt *time.Time `gorm:"default:null" json:"completed_at,omitempty"`
	CreatorID   uuid.UUID  `gorm:"type:uuid;not null" json:"creator_id"`
	ModeratorID *uuid.UUID `gorm:"type:uuid" json:"moderator_id,omitempty"`
}
