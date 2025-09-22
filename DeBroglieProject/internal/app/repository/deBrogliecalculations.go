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
