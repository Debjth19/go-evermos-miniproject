package repository

import (
	"github.com/Debjth19/go-evermos/model"
	
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type AuthRepository interface {
	CheckEmailExists(email string) (bool, error)
	CheckNoTelpExists(noTelp string) (bool, error)
	Register(user model.User) (model.User, error)
	Login(noTelp string) (model.User, error)
}

type authRepository struct {
	db *gorm.DB
}

// NewAuthRepository membuat instance baru dari authRepository
func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{db}
}

// CheckEmailExists mengecek apakah email sudah terdaftar
func (r *authRepository) CheckEmailExists(email string) (bool, error) {
	var user model.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil // Email belum ada (aman)
		}
		return false, err // Error database lainnya
	}
	return true, nil // Email sudah ada
}

// CheckNoTelpExists mengecek apakah no_telp sudah terdaftar
func (r *authRepository) CheckNoTelpExists(noTelp string) (bool, error) {
	var user model.User
	err := r.db.Where("no_telp = ?", noTelp).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil // No Telp belum ada (aman)
		}
		return false, err // Error database lainnya
	}
	return true, nil // No Telp sudah ada
}

// Register menyimpan user baru dan membuat toko baru dalam satu transaksi
func (r *authRepository) Register(user model.User) (model.User, error) {
	// Memulai database transaction
	err := r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Buat User
		if err := tx.Create(&user).Error; err != nil {
			// Jika gagal buat user, batalkan transaksi
			return err
		}

		// 2. Buat Toko
		namaToko := fmt.Sprintf("toko-%s", strings.ToLower(strings.ReplaceAll(user.Nama, " ", "-")))

		toko := model.Toko{
			UserID:   user.ID,
			NamaToko: namaToko,
			UrlFoto:  "", // Default foto kosong
		}

		if err := tx.Create(&toko).Error; err != nil {
			return err
		}

		// Jika semua berhasil, commit transaksi
		return nil
	})

	if err != nil {
		return user, err
	}

	return user, nil
}

// Login mencari user berdasarkan no_telp
func (r *authRepository) Login(noTelp string) (model.User, error) {
	var user model.User
	err := r.db.Where("no_telp = ?", noTelp).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// User tidak ditemukan
			return user, gorm.ErrRecordNotFound
		}
		// Error database lainnya
		return user, err
	}
	
	// User ditemukan
	return user, nil
}