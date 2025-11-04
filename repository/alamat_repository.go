package repository

import (
	"github.com/Debjth19/go-evermos/model"

	"gorm.io/gorm"
)

type AlamatRepository interface {
	Create(alamat model.Alamat) (model.Alamat, error)
	FindAllByUserID(userID uint) ([]model.Alamat, error)
	FindByID(alamatID uint) (model.Alamat, error)
	Update(alamat model.Alamat) (model.Alamat, error)
	Delete(alamatID uint) error
}

type alamatRepository struct {
	db *gorm.DB
}

func NewAlamatRepository(db *gorm.DB) AlamatRepository {
	return &alamatRepository{db}
}

func (r *alamatRepository) Create(alamat model.Alamat) (model.Alamat, error) {
	err := r.db.Create(&alamat).Error
	if err != nil {
		return alamat, err
	}
	return alamat, nil
}

func (r *alamatRepository) FindAllByUserID(userID uint) ([]model.Alamat, error) {
	var alamats []model.Alamat
	err := r.db.Where("user_id = ?", userID).Find(&alamats).Error
	if err != nil {
		return alamats, err
	}
	return alamats, nil
}

func (r *alamatRepository) FindByID(alamatID uint) (model.Alamat, error) {
	var alamat model.Alamat
	err := r.db.Where("id = ?", alamatID).First(&alamat).Error
	if err != nil {
		return alamat, err
	}
	return alamat, nil
}

func (r *alamatRepository) Update(alamat model.Alamat) (model.Alamat, error) {
	err := r.db.Save(&alamat).Error
	if err != nil {
		return alamat, err
	}
	return alamat, nil
}

func (r *alamatRepository) Delete(alamatID uint) error {
	err := r.db.Delete(&model.Alamat{}, alamatID).Error
	if err != nil {
		return err
	}
	return nil
}