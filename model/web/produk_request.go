package web

type ProdukCreateRequest struct {
	NamaProduk    string `validate:"required"`
	CategoryID    uint   `validate:"required"`
	HargaReseler  uint   `validate:"required"`
	HargaKonsumen uint   `validate:"required"`
	Stok          uint   `validate:"required"`
	Deskripsi     string `validate:"required"`
}

type ProdukUpdateRequest struct {
	NamaProduk    string
	CategoryID    uint
	HargaReseler  uint
	HargaKonsumen uint
	Stok          uint
	Deskripsi     string
}