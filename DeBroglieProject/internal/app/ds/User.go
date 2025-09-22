package ds

type User struct {
	ID          uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Email       string `gorm:"type:varchar(25);unique;not null" json:"email"`
	Name        string `gorm:"type:varchar(50);not null" json:"name"`
	Password    string `gorm:"type:varchar(50);not null" json:"-"`
	IsModerator bool   `gorm:"type:boolean;default:false" json:"is_moderator"`
}
