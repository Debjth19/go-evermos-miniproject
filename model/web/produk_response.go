package web

type FotoProdukResponse struct {
	ID        uint   `json:"id"`
	ProductID uint   `json:"product_id"`
	Url       string `json:"url"`
}

type ProdukResponse struct {
	ID            uint                 `json:"id"`
	NamaProduk    string               `json:"nama_produk"`
	Slug          string               `json:"slug"`
	HargaReseler  uint                 `json:"harga_reseler"`
	HargaKonsumen uint                 `json:"harga_konsumen"`
	Stok          uint                 `json:"stok"`
	Deskripsi     string               `json:"deskripsi"`
	Toko          TokoResponse         `json:"toko"`     // Relasi
	Category      KategoriResponse     `json:"category"` // Relasi
	Photos        []FotoProdukResponse `json:"photos"`   // Relasi
}

type PaginatedProdukResponse struct {
	Page  int              `json:"page"`
	Limit int              `json:"limit"`
	Data  []ProdukResponse `json:"data"`
}