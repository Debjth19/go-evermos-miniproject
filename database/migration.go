package database

import (
	"fmt"
	"github.com/Debjth19/go-evermos/config"
	"github.com/Debjth19/go-evermos/model"
)

func MigrateDatabase() {
	err := config.DB.AutoMigrate(
		&model.User{},
		&model.Toko{},
		&model.Alamat{},
		&model.Kategori{},
		&model.Produk{},
		&model.FotoProduk{},
		&model.Transaksi{},
		&model.DetailTransaksi{},
		&model.LogProduk{},
	)
	
	if err != nil {
		panic("Gagal melakukan migrasi database")
	}
	
	fmt.Println("Migrasi database berhasil")
}