package ds

type DeBroglieCalculation struct {
	ID                            uint     `gorm:"primaryKey;autoIncrement"`
	RequestDeBroglieCalculationID uint     `gorm:"not null;uniqueIndex:idx_request_particle"`
	ParticleID                    uint     `gorm:"not null;uniqueIndex:idx_request_particle"`
	Speed                         float64  `gorm:"type:numeric;not null"`
	DeBroglieLength               *float64 `gorm:"type:numeric"`

	Request  RequestDeBroglieCalculation `gorm:"foreignKey:RequestDeBroglieCalculationID"`
	Particle Particle                    `gorm:"foreignKey:ParticleID"`
}
