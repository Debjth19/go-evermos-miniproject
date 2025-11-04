package web

type DetailTransaksiRequest struct {
	ProductID uint `json:"product_id" validate:"required"`
	Kuantitas uint `json:"kuantitas" validate:"required,min=1"`
}

type TransaksiCreateRequest struct {
	MethodBayar  string                   `json:"method_bayar" validate:"required"`
	AlamatKirim  uint                     `json:"alamat_kirim" validate:"required"`
	DetailTrx    []DetailTransaksiRequest `json:"detail_trx" validate:"required,min=1"`
}