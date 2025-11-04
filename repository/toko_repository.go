package repository

import (
	"github.com/Debjth19/go-evermos/helpers"
	"github.com/Debjth19/go-evermos/model"

	"gorm.io/gorm"
)

type TokoRepository interface {
	FindByUserID(userID uint) (model.Toko, error)
	FindByID(tokoID uint) (model.Toko, error)
	Update(toko model.Toko) (model.Toko, error)
	FindAll(pagination helpers.Pagination, search string) ([]model.Toko, error)
}

type tokoRepository struct {
	db *gorm.DB
}

func NewTokoRepository(db *gorm.DB) TokoRepository {
	return &tokoRepository{db}
}

func (r *tokoRepository) FindByUserID(userID uint) (model.Toko, error) {
	var toko model.Toko
	err := r.db.Where("user_id = ?", userID).First(&toko).Error
	if err != nil {
		return toko, err
	}
	return toko, nil
}

func (r *tokoRepository) FindByID(tokoID uint) (model.Toko, error) {
	var toko model.Toko
	err := r.db.Where("id = ?", tokoID).First(&toko).Error
	if err != nil {
		return toko, err
	}
	return toko, nil
}

func (r *tokoRepository) Update(toko model.Toko) (model.Toko, error) {
	err := r.db.Save(&toko).Error
	if err != nil {
		return toko, err
	}
	return toko, nil
}

func (r *tokoRepository) FindAll(pagination helpers.Pagination, search string) ([]model.Toko, error) {
	var tokos []model.Toko
	
	query := r.db.Model(&model.Toko{})

	// Terapkan filter 
	if search != "" {
		query = query.Where("nama_toko LIKE ?", "%"+search+"%")
	}

	// Terapkan pagination 
	offset := (pagination.Page - 1) * pagination.Limit
	err := query.Limit(pagination.Limit).Offset(offset).Find(&tokos).Error
	if err != nil {
		return tokos, err
	}
	
	return tokos, nil
}