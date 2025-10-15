package repository

import (
	"DeBroglieProject/internal/app/ds"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (r *Repository) GetRequestDeBroglieCalculations(status *ds.RequestStatus, startDate, endDate *time.Time, researcherID uuid.UUID, isModerator bool) ([]ds.RequestDeBroglieCalculation, error) {
	var requests []ds.RequestDeBroglieCalculation
	query := r.db.Where("status != ? AND status != ?", ds.RequestStatusDeleted, ds.RequestStatusDraft)

	if !isModerator {
		query = query.Where("researcher_id = ?", researcherID)
	}

	if status != nil {
		query = query.Where("status = ?", *status)
	}

	if startDate != nil {
		query = query.Where("formed_at >= ?", *startDate)
	}

	if endDate != nil {
		query = query.Where("formed_at <= ?", *endDate)
	}

	err := query.Find(&requests).Error
	return requests, err
}

func (r *Repository) GetRequestDeBroglieCalculation(id uint) (ds.RequestDeBroglieCalculation, error) {
	var request ds.RequestDeBroglieCalculation
	err := r.db.Where("id = ? AND status != ?", id, ds.RequestStatusDeleted).First(&request).Error
	return request, err
}

func (r *Repository) GetRequestWithCalculations(id uint) (ds.RequestDeBroglieCalculation, []ds.DeBroglieCalculation, error) {
	req, err := r.GetRequestDeBroglieCalculation(id)
	if err != nil {
		return ds.RequestDeBroglieCalculation{}, nil, err
	}
	var calcs []ds.DeBroglieCalculation
	err = r.db.Preload("Particle").Where("request_de_broglie_calculation_id = ?", id).Find(&calcs).Error
	if err != nil {
		return ds.RequestDeBroglieCalculation{}, nil, err
	}
	return req, calcs, nil
}

func (r *Repository) UpdateRequestDeBroglieCalculation(id uint, request ds.RequestDeBroglieCalculation) error {
	var existingRequest ds.RequestDeBroglieCalculation
	err := r.db.Where("id = ? AND status != ?", id, ds.RequestStatusDeleted).First(&existingRequest).Error
	if err != nil {
		return err
	}

	return r.db.Model(&existingRequest).Updates(request).Error
}

func (r *Repository) UpdateDeBroglieRequestStatus(id uint, newStatus ds.RequestStatus, professorID *uuid.UUID) error {
	var request ds.RequestDeBroglieCalculation
	err := r.db.Where("id = ? AND status != ?", id, ds.RequestStatusDeleted).First(&request).Error
	if err != nil {
		return err
	}

	if !r.isValidStatusTransition(request.Status, newStatus) {
		return fmt.Errorf("недопустимый переход статуса с %s на %s", request.Status, newStatus)
	}

	updates := map[string]interface{}{
		"status": newStatus,
	}

	switch newStatus {
	case ds.RequestStatusFormed:
		updates["formed_at"] = time.Now()
	case ds.RequestStatusCompleted, ds.RequestStatusRejected:
		updates["completed_at"] = time.Now()
		updates["professor_id"] = *professorID
	}

	return r.db.Model(&ds.RequestDeBroglieCalculation{}).Where("id = ?", id).Updates(updates).Error
}

func (r *Repository) isValidStatusTransition(current, new ds.RequestStatus) bool {
	validTransitions := map[ds.RequestStatus][]ds.RequestStatus{
		ds.RequestStatusDraft:     {ds.RequestStatusDeleted, ds.RequestStatusFormed},
		ds.RequestStatusFormed:    {ds.RequestStatusCompleted, ds.RequestStatusRejected},
		ds.RequestStatusCompleted: {},
		ds.RequestStatusRejected:  {},
		ds.RequestStatusDeleted:   {},
	}

	allowedStatuses, exists := validTransitions[current]
	if !exists {
		return false
	}

	for _, status := range allowedStatuses {
		if status == new {
			return true
		}
	}
	return false
}

func (r *Repository) GetDraftRequestDeBroglieCalculationInfo(researcherID uuid.UUID) (ds.RequestDeBroglieCalculation, []ds.DeBroglieCalculation, error) {
	requestDeBroglieCalculation := ds.RequestDeBroglieCalculation{}
	err := r.db.Where("researcher_id = ? AND status = ?", researcherID, ds.RequestStatusDraft).First(&requestDeBroglieCalculation).Error
	if err != nil {
		return ds.RequestDeBroglieCalculation{}, nil, err
	}

	var deBroglieCalculations []ds.DeBroglieCalculation
	err = r.db.Preload("Particle").Where("request_de_broglie_calculation_id = ?", requestDeBroglieCalculation.ID).Find(&deBroglieCalculations).Error
	if err != nil {
		return ds.RequestDeBroglieCalculation{}, nil, err
	}

	return requestDeBroglieCalculation, deBroglieCalculations, nil
}

func (r *Repository) calculateDeBroglieWavelength(mass float64, velocity float64) float64 {
	const planckConstant = 6.62607015e-34
	if mass <= 0 || velocity <= 0 {
		return 0
	}
	return planckConstant / (mass * velocity)
}

func (r *Repository) calculateDeBroglieLengthsForRequest(requestID uint) error {
	var calculations []ds.DeBroglieCalculation
	err := r.db.Preload("Particle").Where("request_de_broglie_calculation_id = ?", requestID).Find(&calculations).Error
	if err != nil {
		return err
	}

	for _, calc := range calculations {
		if calc.Speed == nil {
			continue
		}
		wavelength := r.calculateDeBroglieWavelength(calc.Particle.Mass, *calc.Speed)
		err = r.db.Model(&ds.DeBroglieCalculation{}).Where("id = ?", calc.ID).Update("de_broglie_length", &wavelength).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Repository) AddDeBroglieCalculationToRequest(requestID uint, particleID uint) error {
	deBroglieCalculation := ds.DeBroglieCalculation{
		RequestDeBroglieCalculationID: requestID,
		ParticleID:                    particleID,
		Speed:                         nil,
		DeBroglieLength:               nil,
	}
	err := r.db.Create(&deBroglieCalculation).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) DeleteRequestDeBroglieCalculation(id uint) (int64, error) {
	var existingRequest ds.RequestDeBroglieCalculation
	err := r.db.Where("id = ? AND status != ?", id, ds.RequestStatusDeleted).First(&existingRequest).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, nil
		}
		return 0, err
	}

	result := r.db.Model(&existingRequest).Update("status", ds.RequestStatusDeleted)
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, nil
}

func (r *Repository) FormDeBroglieRequestDraft(id uint, researcherID uuid.UUID) error {
	draft, calcs, err := r.GetDraftRequestDeBroglieCalculationInfo(researcherID)
	if err != nil || draft.ID != id {
		return fmt.Errorf("доступен только черновик текущего пользователя")
	}
	if len(calcs) == 0 {
		return fmt.Errorf("заявка пуста")
	}
	
	if draft.Name == nil {
		return fmt.Errorf("нельзя сформировать заявку без названия")
	}
	
	for _, calc := range calcs {
		if calc.Speed == nil {
			return fmt.Errorf("нельзя сформировать заявку: у частицы %s не указана скорость", calc.Particle.Name)
		}
	}
	
	newStatus := ds.RequestStatusFormed
	return r.UpdateDeBroglieRequestStatus(id, newStatus, nil)
}

func (r *Repository) CompleteDeBroglieRequest(id uint, approve bool, professorID uuid.UUID) error {
	status := ds.RequestStatusRejected
	if approve {
		status = ds.RequestStatusCompleted
	}

	err := r.UpdateDeBroglieRequestStatus(id, status, &professorID)
	if err != nil {
		return err
	}

	if approve && status == ds.RequestStatusCompleted {
		err = r.calculateDeBroglieLengthsForRequest(id)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Repository) CreateRequestDeBroglieCalculationWithParticle(particleID uint, researcherID uuid.UUID) (ds.RequestDeBroglieCalculation, error) {
	requestDeBroglieCalculation := ds.RequestDeBroglieCalculation{
		Name:         nil,
		Status:       ds.RequestStatusDraft,
		ResearcherID: researcherID,
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
