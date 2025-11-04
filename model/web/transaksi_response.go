package web

type DetailTransaksiResponse struct {
	Produk     ProdukResponse `json:"product"` 
	Toko       TokoResponse   `json:"toko"`    
	Kuantitas  uint           `json:"kuantitas"`
	HargaTotal uint           `json:"harga_total"`
}

type TransaksiResponse struct {
	ID            uint                      `json:"id"`
	HargaTotal    uint                      `json:"harga_total"`
	KodeInvoice   string                    `json:"kode_invoice"`
	MethodBayar   string                    `json:"method_bayar"`
	AlamatKirim   AlamatResponse            `json:"alamat_kirim"` 
	DetailTrx     []DetailTransaksiResponse `json:"detail_trx"`
}

type PaginatedTransaksiResponse struct {
	Page  int                 `json:"page"`
	Limit int                 `json:"limit"`
	Data  []TransaksiResponse `json:"data"`
}