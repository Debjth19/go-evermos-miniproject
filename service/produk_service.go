package service

import (
	"github.com/Debjth19/go-evermos/helpers"
	"github.com/Debjth19/go-evermos/model"
	"github.com/Debjth19/go-evermos/model/web"
	"github.com/Debjth19/go-evermos/repository"

	"errors"
	"mime/multipart"
	"strconv"

	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

type ProdukService interface {
	CreateProduk(userID uint, request web.ProdukCreateRequest, files []*multipart.FileHeader) (model.Produk, error)
	GetAllProduk(pagination helpers.Pagination, filterParams map[string]string) ([]model.Produk, error)
	GetProdukByID(produkID uint) (model.Produk, error)
	UpdateProduk(userID uint, produkID uint, request web.ProdukUpdateRequest, files []*multipart.FileHeader) (model.Produk, error)
	DeleteProduk(userID uint, produkID uint) error
}

type produkService struct {
	produkRepository repository.ProdukRepository
	tokoRepository   repository.TokoRepository // Dibutuhkan untuk otorisasi
}

func NewProdukService(produkRepo repository.ProdukRepository, tokoRepo repository.TokoRepository) ProdukService {
	return &produkService{
		produkRepository: produkRepo,
		tokoRepository:   tokoRepo,
	}
}

func (s *produkService) verifyProdukOwnership(userID uint, produkID uint) (model.Produk, error) {
	// 1. Dapatkan toko milik user
	toko, err := s.tokoRepository.FindByUserID(userID)
	if err != nil {
		return model.Produk{}, errors.New("Toko Anda tidak ditemukan")
	}

	// 2. Dapatkan data produk
	produk, err := s.produkRepository.FindByID(produkID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return produk, errors.New("Produk tidak ditemukan")
		}
		return produk, err
	}

	if produk.TokoID != toko.ID {
		return produk, errors.New("Akses ditolak: Anda bukan pemilik produk ini")
	}

	return produk, nil
}

func (s *produkService) CreateProduk(userID uint, request web.ProdukCreateRequest, files []*multipart.FileHeader) (model.Produk, error) {
	// 1. Dapatkan toko milik user
	toko, err := s.tokoRepository.FindByUserID(userID)
	if err != nil {
		return model.Produk{}, errors.New("Toko Anda tidak ditemukan, tidak bisa menambah produk")
	}

	// 2. Simpan file foto (jika ada)
	fotoUrls, err := helpers.SaveUploadedFiles(files, helpers.ProdukImagesPath)
	if err != nil {
		return model.Produk{}, errors.New("Gagal menyimpan foto: " + err.Error())
	}

	// 3. Buat slug
	produkSlug := slug.Make(request.NamaProduk)

	// 4. Buat struct produk
	produk := model.Produk{
		NamaProduk:    request.NamaProduk,
		Slug:          produkSlug,
		HargaReseler:  request.HargaReseler,
		HargaKonsumen: request.HargaKonsumen,
		Stok:          request.Stok,
		Deskripsi:     request.Deskripsi,
		TokoID:        toko.ID, // Set pemilik produk
		CategoryID:    request.CategoryID,
	}

	// 5. Panggil repository untuk create
	newProduk, err := s.produkRepository.Create(produk, fotoUrls)
	if err != nil {
		// Jika create DB gagal, hapus file yang sudah terlanjur di-upload
		helpers.DeleteFiles(fotoUrls, helpers.ProdukImagesPath)
		return newProduk, err
	}

	return newProduk, nil
}

func (s *produkService) parseFilter(filterParams map[string]string) repository.ProdukFilter {
	filter := repository.ProdukFilter{}
	filter.NamaProduk = filterParams["nama_produk"]
	
	if catID, err := strconv.Atoi(filterParams["category_id"]); err == nil {
		filter.CategoryID = uint(catID)
	}
	if tokoID, err := strconv.Atoi(filterParams["toko_id"]); err == nil {
		filter.TokoID = uint(tokoID)
	}
	if minHarga, err := strconv.Atoi(filterParams["min_harga"]); err == nil {
		filter.MinHarga = uint(minHarga)
	}
	if maxHarga, err := strconv.Atoi(filterParams["max_harga"]); err == nil {
		filter.MaxHarga = uint(maxHarga)
	}
	return filter
}

func (s *produkService) GetAllProduk(pagination helpers.Pagination, filterParams map[string]string) ([]model.Produk, error) {
	filter := s.parseFilter(filterParams)
	return s.produkRepository.FindAll(pagination, filter)
}

func (s *produkService) GetProdukByID(produkID uint) (model.Produk, error) {
	produk, err := s.produkRepository.FindByID(produkID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return produk, errors.New("Produk tidak ditemukan")
		}
		return produk, err
	}
	return produk, nil
}

func (s *produkService) UpdateProduk(userID uint, produkID uint, request web.ProdukUpdateRequest, files []*multipart.FileHeader) (model.Produk, error) {
	// 1. Verifikasi kepemilikan
	produk, err := s.verifyProdukOwnership(userID, produkID)
	if err != nil {
		return produk, err
	}

	// 2. Simpan file foto BARU (jika ada)
	var newFotoUrls []string
	if len(files) > 0 {
		// Hapus foto LAMA dari file system
		var oldFotoUrls []string
		for _, foto := range produk.FotoProduk {
			oldFotoUrls = append(oldFotoUrls, foto.Url)
		}
		helpers.DeleteFiles(oldFotoUrls, helpers.ProdukImagesPath)

		// Simpan foto BARU
		newFotoUrls, err = helpers.SaveUploadedFiles(files, helpers.ProdukImagesPath)
		if err != nil {
			return produk, errors.New("Gagal menyimpan foto baru: " + err.Error())
		}
	}

	// 3. Update field
	if request.NamaProduk != "" {
		produk.NamaProduk = request.NamaProduk
		produk.Slug = slug.Make(request.NamaProduk) 
	}
	if request.CategoryID != 0 {
		produk.CategoryID = request.CategoryID
	}
	if request.HargaReseler != 0 {
		produk.HargaReseler = request.HargaReseler
	}
	if request.HargaKonsumen != 0 {
		produk.HargaKonsumen = request.HargaKonsumen
	}
	if request.Stok != 0 {
		produk.Stok = request.Stok
	}
	if request.Deskripsi != "" {
		produk.Deskripsi = request.Deskripsi
	}

	// 4. Update relasi foto di DB (jika ada foto baru)
	if len(newFotoUrls) > 0 {
		// Hapus relasi foto lama di DB
		err = s.produkRepository.DeleteFotosByProductID(produkID, nil) // nil tx
		if err != nil {
			return produk, errors.New("Gagal menghapus relasi foto lama: " + err.Error())
		}
		
		// Buat relasi foto baru di DB
		var newFotos []model.FotoProduk
		for _, url := range newFotoUrls {
			newFotos = append(newFotos, model.FotoProduk{
				ProductID: produkID,
				Url:       url,
			})
		}
		err = s.produkRepository.CreateFotos(newFotos, nil) // nil tx
		if err != nil {
			return produk, errors.New("Gagal membuat relasi foto baru: " + err.Error())
		}
	}
	
	// 5. Simpan perubahan produk
	return s.produkRepository.Update(produk)
}

func (s *produkService) DeleteProduk(userID uint, produkID uint) error {
	// 1. Verifikasi kepemilikan
	produk, err := s.verifyProdukOwnership(userID, produkID)
	if err != nil {
		return err
	}

	// 2. Hapus file foto dari file system
	var fotoUrls []string
	for _, foto := range produk.FotoProduk {
		fotoUrls = append(fotoUrls, foto.Url)
	}
	helpers.DeleteFiles(fotoUrls, helpers.ProdukImagesPath)

	// 3. Hapus produk dari DB (repository akan menangani foto)
	return s.produkRepository.Delete(produkID)
}