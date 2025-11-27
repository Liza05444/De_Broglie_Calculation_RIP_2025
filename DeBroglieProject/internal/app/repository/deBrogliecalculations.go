package repository

import (
	"DeBroglieProject/internal/app/ds"
)

func (r *Repository) RemoveCalculationFromRequest(requestID, particleID uint) (int64, error) {
	result := r.db.Where("request_de_broglie_calculation_id = ? AND particle_id = ?", requestID, particleID).Delete(&ds.DeBroglieCalculation{})
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

func (r *Repository) UpdateCalculationValue(requestID uint, particleID uint, speed float64) error {
	return r.db.Model(&ds.DeBroglieCalculation{}).Where("request_de_broglie_calculation_id = ? AND particle_id = ?", requestID, particleID).Update("speed", speed).Error
}

func (r *Repository) UpdateCalculationDeBroglieLength(calculationID uint, deBroglieLength float64) error {
	return r.db.Model(&ds.DeBroglieCalculation{}).Where("id = ?", calculationID).Update("de_broglie_length", deBroglieLength).Error
}

func (r *Repository) CountCalculationsWithDeBroglieLength(requestID uint) (int64, error) {
	var count int64
	err := r.db.Model(&ds.DeBroglieCalculation{}).Where("request_de_broglie_calculation_id = ? AND de_broglie_length IS NOT NULL", requestID).Count(&count).Error
	return count, err
}

func (r *Repository) CountTotalCalculationsForRequest(requestID uint) (int64, error) {
	var count int64
	err := r.db.Model(&ds.DeBroglieCalculation{}).Where("request_de_broglie_calculation_id = ?", requestID).Count(&count).Error
	return count, err
}

func (r *Repository) GetCalculationByID(calculationID uint) (ds.DeBroglieCalculation, error) {
	var calc ds.DeBroglieCalculation
	err := r.db.Where("id = ?", calculationID).First(&calc).Error
	return calc, err
}
