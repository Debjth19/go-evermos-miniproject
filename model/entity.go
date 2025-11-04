package model

import (
	"time"
)

// User mewakili tabel 'users'
type User struct {
	ID           uint      `gorm:"primaryKey"`
	Nama         string    `gorm:"type:varchar(255)"`
	KataSandi    string    `gorm:"type:varchar(255)"`
	NoTelp       string    `gorm:"type:varchar(20);unique"`
	TanggalLahir time.Time `gorm:"type:date"`
	Pekerjaan    string    `gorm:"type:varchar(100)"`
	Email        string    `gorm:"type:varchar(100);unique"`
	IDProvinsi   string    `gorm:"type:varchar(10)"`
	IDKota       string    `gorm:"type:varchar(10)"`
	Role         string    `gorm:"type:enum('user', 'admin');default:'user'"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Toko         Toko      `gorm:"foreignKey:UserID"` // Relasi one-to-one
	Alamat       []Alamat  `gorm:"foreignKey:UserID"` // Relasi one-to-many
	Transaksi    []Transaksi `gorm:"foreignKey:UserID"` // Relasi one-to-many
}

// Toko mewakili tabel 'toko'
type Toko struct {
	ID        uint      `gorm:"primaryKey"`
	NamaToko  string    `gorm:"type:varchar(100)"`
	UrlFoto   string    `gorm:"type:varchar(255)"`
	UserID    uint      `gorm:"unique"` // Relasi one-to-one
	Produk    []Produk  `gorm:"foreignKey:TokoID"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Alamat mewakili tabel 'alamat'
type Alamat struct {
	ID            uint   `gorm:"primaryKey"`
	JudulAlamat   string `gorm:"type:varchar(100)"`
	NamaPenerima  string `gorm:"type:varchar(100)"`
	NoTelp        string `gorm:"type:varchar(20)"`
	DetailAlamat  string `gorm:"type:text"`
	UserID        uint   // Foreign key ke User
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// Kategori mewakili tabel 'kategori'
type Kategori struct {
	ID            uint      `gorm:"primaryKey"`
	NamaCategory  string    `gorm:"type:varchar(100)"`
	Produk        []Produk  `gorm:"foreignKey:CategoryID"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// Produk mewakili tabel 'produk'
type Produk struct {
	ID             uint      `gorm:"primaryKey"`
	NamaProduk     string    `gorm:"type:varchar(255)"`
	Slug           string    `gorm:"type:varchar(255);unique"`
	HargaReseler   uint
	HargaKonsumen  uint
	Stok           uint
	Deskripsi      string    `gorm:"type:text"`
	TokoID         uint         // Foreign key ke Toko
	CategoryID     uint         // Foreign key ke Kategori
	Toko           Toko         `gorm:"foreignKey:TokoID"`     // Relasi
	Category       Kategori     `gorm:"foreignKey:CategoryID"` // Relasi
	FotoProduk     []FotoProduk `gorm:"foreignKey:ProductID"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// FotoProduk mewakili tabel 'foto_produk'
type FotoProduk struct {
	ID        uint   `gorm:"primaryKey"`
	ProductID uint   // Foreign key ke Produk
	Url       string `gorm:"type:varchar(255)"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Transaksi mewakili tabel 'transaksi'
type Transaksi struct {
	ID              uint   `gorm:"primaryKey"`
	HargaTotal      uint
	KodeInvoice     string `gorm:"type:varchar(100);unique"`
	MethodBayar     string `gorm:"type:varchar(50)"`
	AlamatKirimID   uint   // Foreign key ke Alamat
	UserID          uint   // Foreign key ke User
	Alamat          Alamat `gorm:"foreignKey:AlamatKirimID"` // Relasi
	DetailTransaksi []DetailTransaksi `gorm:"foreignKey:TransaksiID"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// DetailTransaksi mewakili tabel 'detail_transaksi'
type DetailTransaksi struct {
	ID          uint   `gorm:"primaryKey"`
	TransaksiID uint   // Foreign key ke Transaksi
	ProductID   uint   // Foreign key ke Produk
	TokoID      uint   // Foreign key ke Toko
	Kuantitas   uint
	HargaTotal  uint
	Produk      Produk `gorm:"foreignKey:ProductID"` // Relasi
	Toko        Toko   `gorm:"foreignKey:TokoID"`    // Relasi
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// LogProduk mewakili tabel 'log_produk'
type LogProduk struct {
	ID            uint   `gorm:"primaryKey"`
	TransaksiID   uint   // Foreign key ke Transaksi
	ProductID     uint   // ID produk asli
	NamaProduk    string `gorm:"type:varchar(255)"` // Snapshot data
	Slug          string `gorm:"type:varchar(255)"` // Snapshot data
	HargaReseler  uint   // Snapshot data
	HargaKonsumen uint   // Snapshot data
	Deskripsi     string `gorm:"type:text"`         // Snapshot data
	TokoID        uint   // Snapshot data
	CategoryID    uint   // Snapshot data
	CreatedAt     time.Time
	UpdatedAt     time.Time
}