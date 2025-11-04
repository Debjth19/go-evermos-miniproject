package service

import (
	"github.com/Debjth19/go-evermos/model"
	"github.com/Debjth19/go-evermos/model/web"
	"github.com/Debjth19/go-evermos/repository"

	"errors"
	"gorm.io/gorm"
)

type KategoriService interface {
	CreateKategori(request web.KategoriCreateRequest) (model.Kategori, error)
	GetAllKategori() ([]model.Kategori, error)
	GetKategoriByID(kategoriID uint) (model.Kategori, error)
	UpdateKategori(kategoriID uint, request web.KategoriUpdateRequest) (model.Kategori, error)
	DeleteKategori(kategoriID uint) error
}

type kategoriService struct {
	kategoriRepository repository.KategoriRepository
}

func NewKategoriService(kategoriRepo repository.KategoriRepository) KategoriService {
	return &kategoriService{kategoriRepository: kategoriRepo}
}

func (s *kategoriService) CreateKategori(request web.KategoriCreateRequest) (model.Kategori, error) {
	// Cek duplikat nama
	exists, err := s.kategoriRepository.CheckCategoryExists(request.NamaCategory)
	if err != nil {
		return model.Kategori{}, err
	}
	if exists {
		return model.Kategori{}, errors.New("Nama kategori sudah ada")
	}

	kategori := model.Kategori{
		NamaCategory: request.NamaCategory,
	}

	return s.kategoriRepository.Create(kategori)
}

func (s *kategoriService) GetAllKategori() ([]model.Kategori, error) {
	return s.kategoriRepository.FindAll()
}

func (s *kategoriService) GetKategoriByID(kategoriID uint) (model.Kategori, error) {
	kategori, err := s.kategoriRepository.FindByID(kategoriID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return kategori, errors.New("Kategori tidak ditemukan")
		}
		return kategori, err
	}
	return kategori, nil
}

func (s *kategoriService) UpdateKategori(kategoriID uint, request web.KategoriUpdateRequest) (model.Kategori, error) {
	// 1. Cek apakah kategori ada
	kategori, err := s.GetKategoriByID(kategoriID)
	if err != nil {
		return kategori, err // Mengembalikan error "Kategori tidak ditemukan"
	}

	// 2. Cek duplikat nama (jika nama diubah)
	if request.NamaCategory != kategori.NamaCategory {
		exists, err := s.kategoriRepository.CheckCategoryExists(request.NamaCategory)
		if err != nil {
			return kategori, err
		}
		if exists {
			return kategori, errors.New("Nama kategori sudah ada")
		}
	}

	// 3. Update data
	kategori.NamaCategory = request.NamaCategory
	return s.kategoriRepository.Update(kategori)
}

func (s *kategoriService) DeleteKategori(kategoriID uint) error {
	// 1. Cek apakah kategori ada
	_, err := s.GetKategoriByID(kategoriID)
	if err != nil {
		return err // Mengembalikan error "Kategori tidak ditemukan"
	}
	
	// 2. Hapus
	return s.kategoriRepository.Delete(kategoriID)
}