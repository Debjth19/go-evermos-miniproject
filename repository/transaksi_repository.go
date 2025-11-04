package repository

import (
	"github.com/Debjth19/go-evermos/model"

	"gorm.io/gorm"
)

type TransaksiRepository interface {
	Create(tx *gorm.DB, transaksi *model.Transaksi) error
	CreateDetail(tx *gorm.DB, details []model.DetailTransaksi) error
	CreateLog(tx *gorm.DB, logs []model.LogProduk) error
	FindMyTransactions(userID uint) ([]model.Transaksi, error)
	FindMyTransactionByID(userID, trxID uint) (model.Transaksi, error)
}

type transaksiRepository struct {
	db *gorm.DB
}

func NewTransaksiRepository(db *gorm.DB) TransaksiRepository {
	return &transaksiRepository{db}
}

// Create menyimpan transaksi utama
func (r *transaksiRepository) Create(tx *gorm.DB, transaksi *model.Transaksi) error {
	return tx.Create(transaksi).Error
}

// CreateDetail menyimpan item-item detail
func (r *transaksiRepository) CreateDetail(tx *gorm.DB, details []model.DetailTransaksi) error {
	return tx.Create(&details).Error
}

// CreateLog menyimpan snapshot produk
func (r *transaksiRepository) CreateLog(tx *gorm.DB, logs []model.LogProduk) error {
	return tx.Create(&logs).Error
}

// preloads adalah helper untuk query GET agar data relasinya ikut terambil
func (r *transaksiRepository) preloads() *gorm.DB {
	return r.db.
		Preload("Alamat"). // Nama relasi di model.Transaksi
		Preload("DetailTransaksi"). // Nama relasi
		Preload("DetailTransaksi.Produk"). // Relasi di dalam DetailTransaksi
		Preload("DetailTransaksi.Produk.Toko").
		Preload("DetailTransaksi.Produk.Category").
		Preload("DetailTransaksi.Produk.FotoProduk").
		Preload("DetailTransaksi.Toko") // Relasi Toko di DetailTransaksi
}

// FindMyTransactions mengambil semua transaksi milik user
func (r *transaksiRepository) FindMyTransactions(userID uint) ([]model.Transaksi, error) {
	var transaksis []model.Transaksi
	err := r.preloads().
		Where("user_id = ?", userID).
		Order("created_at desc"). // Tampilkan yang terbaru dulu
		Find(&transaksis).Error
	return transaksis, err
}

// FindMyTransactionByID mengambil satu transaksi milik user
func (r *transaksiRepository) FindMyTransactionByID(userID, trxID uint) (model.Transaksi, error) {
	var transaksi model.Transaksi
	err := r.preloads().
		Where("user_id = ? AND id = ?", userID, trxID).
		First(&transaksi).Error
	return transaksi, err
}