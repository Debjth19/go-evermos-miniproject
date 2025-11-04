package service

import (
	"github.com/Debjth19/go-evermos/model"
	"github.com/Debjth19/go-evermos/model/web"
	"github.com/Debjth19/go-evermos/repository"

	"errors"
)

type AlamatService interface {
	CreateAlamat(userID uint, request web.AlamatCreateRequest) (model.Alamat, error)
	GetAllAlamat(userID uint) ([]model.Alamat, error)
	GetAlamatByID(userID uint, alamatID uint) (model.Alamat, error)
	UpdateAlamat(userID uint, alamatID uint, request web.AlamatUpdateRequest) (model.Alamat, error)
	DeleteAlamat(userID uint, alamatID uint) error
}

type alamatService struct {
	alamatRepository repository.AlamatRepository
}

func NewAlamatService(alamatRepo repository.AlamatRepository) AlamatService {
	return &alamatService{alamatRepository: alamatRepo}
}

// verifyAlamatOwnership adalah fungsi helper internal untuk mengecek kepemilikan
func (s *alamatService) verifyAlamatOwnership(userID uint, alamatID uint) (model.Alamat, error) {
	alamat, err := s.alamatRepository.FindByID(alamatID)
	if err != nil {
		return alamat, errors.New("Alamat tidak ditemukan")
	}

	if alamat.UserID != userID {
		return alamat, errors.New("Akses ditolak: Anda bukan pemilik alamat ini")
	}

	return alamat, nil
}

func (s *alamatService) CreateAlamat(userID uint, request web.AlamatCreateRequest) (model.Alamat, error) {
	// Buat struct model Alamat baru
	alamat := model.Alamat{
		JudulAlamat:  request.JudulAlamat,
		NamaPenerima: request.NamaPenerima,
		NoTelp:       request.NoTelp,
		DetailAlamat: request.DetailAlamat,
		UserID:       userID, // Set pemilik alamat
	}

	// Simpan ke database
	newAlamat, err := s.alamatRepository.Create(alamat)
	if err != nil {
		return newAlamat, err
	}

	return newAlamat, nil
}

func (s *alamatService) GetAllAlamat(userID uint) ([]model.Alamat, error) {
	alamats, err := s.alamatRepository.FindAllByUserID(userID)
	if err != nil {
		return alamats, err
	}
	return alamats, nil
}

func (s *alamatService) GetAlamatByID(userID uint, alamatID uint) (model.Alamat, error) {
	// Verifikasi kepemilikan sebelum mengembalikan data
	alamat, err := s.verifyAlamatOwnership(userID, alamatID)
	if err != nil {
		return alamat, err
	}
	
	return alamat, nil
}

func (s *alamatService) UpdateAlamat(userID uint, alamatID uint, request web.AlamatUpdateRequest) (model.Alamat, error) {
	// 1. Verifikasi kepemilikan
	alamat, err := s.verifyAlamatOwnership(userID, alamatID)
	if err != nil {
		return alamat, err
	}

	// 2. Update field jika diisi
	if request.JudulAlamat != "" {
		alamat.JudulAlamat = request.JudulAlamat
	}
	if request.NamaPenerima != "" {
		alamat.NamaPenerima = request.NamaPenerima
	}
	if request.NoTelp != "" {
		alamat.NoTelp = request.NoTelp
	}
	if request.DetailAlamat != "" {
		alamat.DetailAlamat = request.DetailAlamat
	}

	// 3. Simpan perubahan
	updatedAlamat, err := s.alamatRepository.Update(alamat)
	if err != nil {
		return updatedAlamat, err
	}

	return updatedAlamat, nil
}

func (s *alamatService) DeleteAlamat(userID uint, alamatID uint) error {
	// 1. Verifikasi kepemilikan
	_, err := s.verifyAlamatOwnership(userID, alamatID)
	if err != nil {
		return err
	}

	// 2. Hapus alamat
	err = s.alamatRepository.Delete(alamatID)
	if err != nil {
		return err
	}

	return nil
}