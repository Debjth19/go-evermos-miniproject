package repository

import (
	"github.com/Debjth19/go-evermos/helpers"
	"github.com/Debjth19/go-evermos/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Struct untuk filter
type ProdukFilter struct {
	NamaProduk string
	CategoryID uint
	TokoID     uint
	MinHarga   uint
	MaxHarga   uint
}

type ProdukRepository interface {
	Create(produk model.Produk, fotoUrls []string) (model.Produk, error)
	FindAll(pagination helpers.Pagination, filter ProdukFilter) ([]model.Produk, error)
	FindByID(produkID uint) (model.Produk, error)
	Update(produk model.Produk) (model.Produk, error)
	Delete(produkID uint) error
	DeleteFotosByProductID(produkID uint, tx *gorm.DB) error
	CreateFotos(fotos []model.FotoProduk, tx *gorm.DB) error

	FindByIDForUpdate(tx *gorm.DB, produkID uint) (model.Produk, error)
    UpdateStok(tx *gorm.DB, produkID uint, newStok uint) error
}

type produkRepository struct {
	db *gorm.DB
}

func NewProdukRepository(db *gorm.DB) ProdukRepository {
	return &produkRepository{db}
}

// Create membuat produk dan foto-fotonya dalam satu transaksi
func (r *produkRepository) Create(produk model.Produk, fotoUrls []string) (model.Produk, error) {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Buat Produk
		if err := tx.Create(&produk).Error; err != nil {
			return err
		}

		// 2. Buat FotoProduk
		if len(fotoUrls) > 0 {
			var fotos []model.FotoProduk
			for _, url := range fotoUrls {
				fotos = append(fotos, model.FotoProduk{
					ProductID: produk.ID,
					Url:       url,
				})
			}
			if err := tx.Create(&fotos).Error; err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return produk, err
	}
	return produk, nil
}

// FindAll mengambil semua produk dengan pagination dan filter
func (r *produkRepository) FindAll(pagination helpers.Pagination, filter ProdukFilter) ([]model.Produk, error) {
	var produks []model.Produk

	// Siapkan query
	query := r.db.Model(&model.Produk{}).
		Preload("Toko").
		Preload("Category").
		Preload("FotoProduk")

	// Terapkan filter
	if filter.NamaProduk != "" {
		query = query.Where("nama_produk LIKE ?", "%"+filter.NamaProduk+"%")
	}
	if filter.CategoryID != 0 {
		query = query.Where("category_id = ?", filter.CategoryID)
	}
	if filter.TokoID != 0 {
		query = query.Where("toko_id = ?", filter.TokoID)
	}
	if filter.MinHarga != 0 {
		query = query.Where("harga_konsumen >= ?", filter.MinHarga)
	}
	if filter.MaxHarga != 0 {
		query = query.Where("harga_konsumen <= ?", filter.MaxHarga)
	}

	// Terapkan pagination 
	offset := (pagination.Page - 1) * pagination.Limit
	err := query.Limit(pagination.Limit).Offset(offset).Find(&produks).Error
	if err != nil {
		return produks, err
	}
	
	return produks, nil
}

// FindByID mengambil produk tunggal
func (r *produkRepository) FindByID(produkID uint) (model.Produk, error) {
	var produk model.Produk
	err := r.db.
		Preload("Toko").
		Preload("Category").
		Preload("FotoProduk").
		Where("id = ?", produkID).First(&produk).Error
	return produk, err
}

// Update menyimpan perubahan pada produk
func (r *produkRepository) Update(produk model.Produk) (model.Produk, error) {
	err := r.db.Save(&produk).Error
	return produk, err
}

// DeleteFotosByProductID menghapus semua foto terkait produk 
func (r *produkRepository) DeleteFotosByProductID(produkID uint, tx *gorm.DB) error {
	return tx.Where("product_id = ?", produkID).Delete(&model.FotoProduk{}).Error
}

// CreateFotos menambahkan foto baru 
func (r *produkRepository) CreateFotos(fotos []model.FotoProduk, tx *gorm.DB) error {
	return tx.Create(&fotos).Error
}

// Delete menghapus produk
func (r *produkRepository) Delete(produkID uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Hapus FotoProduk
		if err := tx.Where("product_id = ?", produkID).Delete(&model.FotoProduk{}).Error; err != nil {
			return err
		}
		// 2. Hapus Produk
		if err := tx.Delete(&model.Produk{}, produkID).Error; err != nil {
			return err
		}
		return nil
	})
}

// FindByIDForUpdate mengambil produk dan mengunci barisnya
func (r *produkRepository) FindByIDForUpdate(tx *gorm.DB, produkID uint) (model.Produk, error) {
	var produk model.Produk
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ?", produkID).First(&produk).Error
	return produk, err
}

// UpdateStok hanya memperbarui stok
func (r *produkRepository) UpdateStok(tx *gorm.DB, produkID uint, newStok uint) error {
	return tx.Model(&model.Produk{}).Where("id = ?", produkID).Update("stok", newStok).Error
}