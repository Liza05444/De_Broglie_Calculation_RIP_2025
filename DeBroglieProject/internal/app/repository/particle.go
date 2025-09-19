package repository

import (
	"fmt"

	"DeBroglieProject/internal/app/ds"
)

func (r *Repository) GetParticles() ([]ds.Particle, error) {
	var particles []ds.Particle
	err := r.db.Where("is_deleted = ?", false).Find(&particles).Error
	if err != nil {
		return nil, err
	}
	if len(particles) == 0 {
		return nil, fmt.Errorf("частиц нет")
	}
	return particles, nil
}

func (r *Repository) GetParticle(id int) (ds.Particle, error) {
	particle := ds.Particle{}
	err := r.db.Where("id = ? AND is_deleted = ?", id, false).First(&particle).Error
	if err != nil {
		return ds.Particle{}, err
	}
	return particle, nil
}

func (r *Repository) GetParticlesByName(name string) ([]ds.Particle, error) {
	var particles []ds.Particle
	err := r.db.Where("name ILIKE ? AND is_deleted = ?", "%"+name+"%", false).Find(&particles).Error
	if err != nil {
		return nil, err
	}
	return particles, nil
}
