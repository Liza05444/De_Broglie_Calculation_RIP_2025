package ds

type User struct {
	ID          uint   `gorm:"primaryKey;autoIncrement"`
	Email       string `gorm:"type:varchar(25);unique;not null"`
	Name        string `gorm:"type:varchar(50);not null"`
	Password    string `gorm:"type:varchar(50);not null"`
	IsModerator bool   `gorm:"type:boolean;default:false"`
}
