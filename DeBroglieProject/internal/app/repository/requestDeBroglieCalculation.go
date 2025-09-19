package repository

import (
	"DeBroglieProject/internal/app/ds"
	"database/sql"
)

func (r *Repository) GetDraftRequestDeBroglieCalculationInfo() (ds.RequestDeBroglieCalculation, []ds.DeBroglieCalculation, error) {
	creatorID := 1

	requestDeBroglieCalculation := ds.RequestDeBroglieCalculation{}
	err := r.db.Where("creator_id = ? AND status = ?", creatorID, ds.RequestStatusDraft).First(&requestDeBroglieCalculation).Error
	if err != nil {
		return ds.RequestDeBroglieCalculation{}, nil, err
	}

	var deBroglieCalculations []ds.DeBroglieCalculation
	err = r.db.Where("request_de_broglie_calculation_id = ?", requestDeBroglieCalculation.ID).Find(&deBroglieCalculations).Error
	if err != nil {
		return ds.RequestDeBroglieCalculation{}, nil, err
	}

	return requestDeBroglieCalculation, deBroglieCalculations, nil
}

func (r *Repository) CreateRequestDeBroglieCalculation(particleID uint) (ds.RequestDeBroglieCalculation, error) {
	requestDeBroglieCalculation := ds.RequestDeBroglieCalculation{
		Name:      "Эксперимент",
		Status:    ds.RequestStatusDraft,
		CreatorID: 1,
	}
	err := r.db.Create(&requestDeBroglieCalculation).Error
	if err != nil {
		return ds.RequestDeBroglieCalculation{}, err
	}

	err = r.AddDeBroglieCalculationToRequest(requestDeBroglieCalculation.ID, particleID)
	if err != nil {
		return ds.RequestDeBroglieCalculation{}, err
	}

	return requestDeBroglieCalculation, nil
}

func (r *Repository) calculateDeBroglieWavelength(mass float64, velocity float64) float64 {
	const planckConstant = 6.62607015e-34
	if mass <= 0 || velocity <= 0 {
		return 0
	}
	return planckConstant / (mass * velocity)
}

func (r *Repository) AddDeBroglieCalculationToRequest(requestID uint, particleID uint) error {
	var particle ds.Particle
	err := r.db.Where("id = ?", particleID).First(&particle).Error
	if err != nil {
		return err
	}

	speed := 1000.0
	wavelength := r.calculateDeBroglieWavelength(particle.Mass, speed)

	deBroglieCalculation := ds.DeBroglieCalculation{
		RequestDeBroglieCalculationID: requestID,
		ParticleID:                    particleID,
		Speed:                         speed,
		DeBroglieLength:               sql.NullFloat64{Float64: wavelength, Valid: true},
	}
	err = r.db.Create(&deBroglieCalculation).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) DeleteRequestDeBroglieCalculation(requestID int) error {
	result := r.db.Exec("UPDATE request_de_broglie_calculations SET status = ? WHERE id = ?", ds.RequestStatusDeleted, requestID)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
