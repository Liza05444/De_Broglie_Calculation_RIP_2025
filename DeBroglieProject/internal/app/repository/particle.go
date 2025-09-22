package repository

import (
	"DeBroglieProject/internal/app/ds"
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
)

func (r *Repository) GetParticles(particleName string) ([]ds.Particle, error) {
	var particles []ds.Particle
	query := r.db.Where("is_deleted = ?", false)

	if particleName != "" {
		query = query.Where("name ILIKE ?", "%"+particleName+"%")
	}

	err := query.Find(&particles).Error
	return particles, err
}

func (r *Repository) GetParticle(id uint) (ds.Particle, error) {
	var particle ds.Particle
	err := r.db.Where("id = ? AND is_deleted = ?", id, false).First(&particle).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ds.Particle{}, errors.New("particle not found")
		}
		return ds.Particle{}, err
	}
	return particle, nil
}

func (r *Repository) CreateParticle(particle ds.Particle) (ds.Particle, error) {
	err := r.db.Create(&particle).Error
	return particle, err
}

func (r *Repository) UpdateParticle(id uint, particle ds.Particle) error {
	tx := r.db.Model(&ds.Particle{}).Where("id = ? AND is_deleted = ?", id, false).Updates(particle)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return errors.New("particle not found")
	}
	return nil
}

func (r *Repository) DeleteParticle(id uint) error {
	tx := r.db.Model(&ds.Particle{}).Where("id = ?", id).Update("is_deleted", true)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return errors.New("particle not found")
	}
	return nil
}

func (r *Repository) UpdateParticleImage(id uint, imagePath string) error {
	tx := r.db.Model(&ds.Particle{}).Where("id = ? AND is_deleted = ?", id, false).Update("image", imagePath)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return errors.New("particle not found")
	}
	return nil
}

func (r *Repository) UploadFileToMinIO(ctx context.Context, fileName string, fileReader io.Reader, fileSize int64, contentType string) error {
	// Проверяем доступность MinIO
	_, err := r.minio.ListBuckets(ctx)
	if err != nil {
		return fmt.Errorf("MinIO is not accessible: %v", err)
	}

	// Загружаем файл
	_, err = r.minio.PutObject(ctx, r.bucket, fileName, fileReader, fileSize, minio.PutObjectOptions{
		ContentType: contentType,
	})
	return err
}

func (r *Repository) DeleteFileFromMinIO(ctx context.Context, fileName string) error {
	err := r.minio.RemoveObject(ctx, r.bucket, fileName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete file from MinIO: %w", err)
	}
	return nil
}
