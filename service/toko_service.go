package service

import (
	"github.com/Debjth19/go-evermos/helpers"
	"github.com/Debjth19/go-evermos/model"
	"github.com/Debjth19/go-evermos/model/web"
	"github.com/Debjth19/go-evermos/repository"

	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"time"

	"gorm.io/gorm"
)

type TokoService interface {
	GetMyToko(userID uint) (model.Toko, error)
	GetTokoByID(tokoID uint) (model.Toko, error)
	GetAllToko(pagination helpers.Pagination, search string) ([]model.Toko, error)
	UpdateToko(userID uint, tokoID uint, request web.TokoUpdateRequest, file *multipart.FileHeader) (model.Toko, error)
}

type tokoService struct {
	tokoRepository repository.TokoRepository
}

func NewTokoService(tokoRepo repository.TokoRepository) TokoService {
	return &tokoService{tokoRepository: tokoRepo}
}

// verifyTokoOwnership adalah helper internal untuk mengecek kepemilikan toko
func (s *tokoService) verifyTokoOwnership(userID uint, tokoID uint) (model.Toko, error) {
	toko, err := s.tokoRepository.FindByID(tokoID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return toko, errors.New("Toko tidak ditemukan")
		}
		return toko, err
	}

	if toko.UserID != userID {
		return toko, errors.New("Akses ditolak: Anda bukan pemilik toko ini")
	}

	return toko, nil
}

// GetMyToko mengambil data toko milik user yang sedang login
func (s *tokoService) GetMyToko(userID uint) (model.Toko, error) {
	toko, err := s.tokoRepository.FindByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return toko, errors.New("Toko tidak ditemukan")
		}
		return toko, err
	}
	return toko, nil
}

// GetTokoByID mengambil data toko (publik)
func (s *tokoService) GetTokoByID(tokoID uint) (model.Toko, error) {
	toko, err := s.tokoRepository.FindByID(tokoID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return toko, errors.New("Toko tidak ditemukan")
		}
		return toko, err
	}
	return toko, nil
}

// GetAllToko mengambil semua toko dengan pagination dan filter
func (s *tokoService) GetAllToko(pagination helpers.Pagination, search string) ([]model.Toko, error) {
	tokos, err := s.tokoRepository.FindAll(pagination, search)
	if err != nil {
		return tokos, err
	}
	return tokos, nil
}

// UpdateToko menangani logika update (nama & upload foto)
func (s *tokoService) UpdateToko(userID uint, tokoID uint, request web.TokoUpdateRequest, file *multipart.FileHeader) (model.Toko, error) {
	// 1. Verifikasi kepemilikan toko
	toko, err := s.verifyTokoOwnership(userID, tokoID)
	if err != nil {
		return toko, err // Error "Akses ditolak" atau "Tidak ditemukan"
	}

	// 2. Proses Upload File (jika ada file baru)
	if file != nil {
		// Hapus foto lama jika ada
		if toko.UrlFoto != "" {
			oldPath := fmt.Sprintf("./public/images/toko/%s", toko.UrlFoto)
			_ = os.Remove(oldPath) // Abaikan error jika file tidak ada
		}

		// Buat nama file unik
		filename := fmt.Sprintf("%d-%s", time.Now().UnixNano(), file.Filename)
		filePath := fmt.Sprintf("./public/images/toko/%s", filename)

		// Buka file yang diupload
		src, err := file.Open()
		if err != nil {
			return toko, err
		}
		defer src.Close()

		// Buat file baru di server
		dst, err := os.Create(filePath)
		if err != nil {
			return toko, err
		}
		defer dst.Close()

		// Salin file
		if _, err = io.Copy(dst, src); err != nil {
			return toko, err
		}

		// Simpan nama file baru ke struct
		toko.UrlFoto = filename
	}

	// 3. Update nama toko (jika diisi)
	if request.NamaToko != "" {
		toko.NamaToko = request.NamaToko
	}

	// 4. Simpan perubahan ke database
	updatedToko, err := s.tokoRepository.Update(toko)
	if err != nil {
		return updatedToko, err
	}

	return updatedToko, nil
}