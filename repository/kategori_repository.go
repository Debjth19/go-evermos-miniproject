package repository

import (
	"github.com/Debjth19/go-evermos/model"

	"gorm.io/gorm"
)

type KategoriRepository interface {
	Create(kategori model.Kategori) (model.Kategori, error)
	FindAll() ([]model.Kategori, error)
	FindByID(kategoriID uint) (model.Kategori, error)
	Update(kategori model.Kategori) (model.Kategori, error)
	Delete(kategoriID uint) error
	CheckCategoryExists(nama string) (bool, error)
}

type kategoriRepository struct {
	db *gorm.DB
}

func NewKategoriRepository(db *gorm.DB) KategoriRepository {
	return &kategoriRepository{db}
}

func (r *kategoriRepository) Create(kategori model.Kategori) (model.Kategori, error) {
	err := r.db.Create(&kategori).Error
	return kategori, err
}

func (r *kategoriRepository) FindAll() ([]model.Kategori, error) {
	var kategoris []model.Kategori
	err := r.db.Find(&kategoris).Error
	return kategoris, err
}

func (r *kategoriRepository) FindByID(kategoriID uint) (model.Kategori, error) {
	var kategori model.Kategori
	err := r.db.Where("id = ?", kategoriID).First(&kategori).Error
	return kategori, err
}

func (r *kategoriRepository) Update(kategori model.Kategori) (model.Kategori, error) {
	err := r.db.Save(&kategori).Error
	return kategori, err
}

func (r *kategoriRepository) Delete(kategoriID uint) error {
	err := r.db.Delete(&model.Kategori{}, kategoriID).Error
	return err
}

// CheckCategoryExists mengecek apakah nama kategori sudah ada
func (r *kategoriRepository) CheckCategoryExists(nama string) (bool, error) {
	var kategori model.Kategori
	err := r.db.Where("nama_category = ?", nama).First(&kategori).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil // Belum ada (aman)
		}
		return false, err // Error lain
	}
	return true, nil // Sudah ada
}