package service

import (
	"github.com/Debjth19/go-evermos/model"
	"github.com/Debjth19/go-evermos/model/web"
	"github.com/Debjth19/go-evermos/repository"

	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type TransaksiService interface {
	CreateTransaksi(userID uint, request web.TransaksiCreateRequest) (model.Transaksi, error)
	GetMyTransactions(userID uint) ([]model.Transaksi, error)
	GetMyTransactionByID(userID uint, trxID uint) (model.Transaksi, error)
}

type transaksiService struct {
	db                  *gorm.DB // Dibutuhkan untuk memulai transaction
	transaksiRepository repository.TransaksiRepository
	produkRepository    repository.ProdukRepository // Dibutuhkan untuk cek stok & update
	alamatRepository    repository.AlamatRepository // Dibutuhkan untuk cek kepemilikan alamat
}

func NewTransaksiService(db *gorm.DB, trxRepo repository.TransaksiRepository, produkRepo repository.ProdukRepository, alamatRepo repository.AlamatRepository) TransaksiService {
	return &transaksiService{
		db:                  db,
		transaksiRepository: trxRepo,
		produkRepository:    produkRepo,
		alamatRepository:    alamatRepo,
	}
}

// CreateTransaksi menangani semua logika pembuatan transaksi
func (s *transaksiService) CreateTransaksi(userID uint, request web.TransaksiCreateRequest) (model.Transaksi, error) {
	var transaksi model.Transaksi
	var hargaTotalTransaksi uint = 0

	err := s.db.Transaction(func(tx *gorm.DB) error {
		// 1. Verifikasi Alamat Kirim
		alamat, err := s.alamatRepository.FindByID(request.AlamatKirim)
		if err != nil {
			return errors.New("Alamat kirim tidak ditemukan")
		}
		// Cek kepemilikan alamat
		if alamat.UserID != userID {
			return errors.New("Akses ditolak: Alamat kirim bukan milik Anda")
		}

		// Siapkan slice untuk detail dan log
		var details []model.DetailTransaksi
		var logs []model.LogProduk

		// 2. Loop setiap item produk di keranjang
		for _, item := range request.DetailTrx {
			produk, err := s.produkRepository.FindByIDForUpdate(tx, item.ProductID)
			if err != nil {
				return fmt.Errorf("Produk dengan ID %d tidak ditemukan", item.ProductID)
			}

			// Cek Stok
			if produk.Stok < item.Kuantitas {
				return fmt.Errorf("Stok tidak mencukupi untuk produk: %s", produk.NamaProduk)
			}

			// Hitung harga total untuk item ini
			hargaTotalItem := produk.HargaKonsumen * item.Kuantitas
			hargaTotalTransaksi += hargaTotalItem

			// Kurangi stok
			newStok := produk.Stok - item.Kuantitas
			if err := s.produkRepository.UpdateStok(tx, produk.ID, newStok); err != nil {
				return fmt.Errorf("Gagal update stok untuk: %s", produk.NamaProduk)
			}

			// Siapkan data DetailTransaksi
			details = append(details, model.DetailTransaksi{
				// TransaksiID akan diisi nanti setelah Transaksi dibuat
				ProductID:   produk.ID,
				TokoID:      produk.TokoID,
				Kuantitas:   item.Kuantitas,
				HargaTotal:  hargaTotalItem,
			})

			logs = append(logs, model.LogProduk{
				// TransaksiID akan diisi nanti
				ProductID:     produk.ID,
				NamaProduk:    produk.NamaProduk,
				Slug:          produk.Slug,
				HargaReseler:  produk.HargaReseler,
				HargaKonsumen: produk.HargaKonsumen,
				Deskripsi:     produk.Deskripsi,
				TokoID:        produk.TokoID,
				CategoryID:    produk.CategoryID,
			})
		}

		// 3. Buat Transaksi utama
		kodeInvoice := fmt.Sprintf("INV-%d-%d", userID, time.Now().Unix())
		transaksi = model.Transaksi{
			HargaTotal:    hargaTotalTransaksi,
			KodeInvoice:   kodeInvoice,
			MethodBayar:   request.MethodBayar,
			AlamatKirimID: request.AlamatKirim,
			UserID:        userID,
		}
		
		if err := s.transaksiRepository.Create(tx, &transaksi); err != nil {
			return errors.New("Gagal membuat transaksi")
		}

		// 4. Update TransaksiID di details dan logs
		for i := range details {
			details[i].TransaksiID = transaksi.ID
		}
		for i := range logs {
			logs[i].TransaksiID = transaksi.ID
		}

		// 5. Simpan DetailTransaksi
		if err := s.transaksiRepository.CreateDetail(tx, details); err != nil {
			return errors.New("Gagal menyimpan detail transaksi")
		}

		// 6. Simpan LogProduk
		if err := s.transaksiRepository.CreateLog(tx, logs); err != nil {
			return errors.New("Gagal menyimpan log produk")
		}

		return nil
	}) 

	if err != nil {
		return model.Transaksi{}, err
	}

	// Kembalikan ID transaksi yang baru dibuat
	return transaksi, nil
}

// GetMyTransactions mengambil semua transaksi milik user
func (s *transaksiService) GetMyTransactions(userID uint) ([]model.Transaksi, error) {
	transaksis, err := s.transaksiRepository.FindMyTransactions(userID)
	if err != nil {
		return nil, err
	}
	return transaksis, nil
}

// GetMyTransactionByID mengambil satu transaksi milik user (Ketentuan No. 15)
func (s *transaksiService) GetMyTransactionByID(userID uint, trxID uint) (model.Transaksi, error) {
	transaksi, err := s.transaksiRepository.FindMyTransactionByID(userID, trxID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return transaksi, errors.New("Transaksi tidak ditemukan atau bukan milik Anda")
		}
		return transaksi, err
	}
	return transaksi, nil
}