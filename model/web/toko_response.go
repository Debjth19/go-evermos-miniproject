package web

type TokoResponse struct {
	ID       uint   `json:"id"`
	NamaToko string `json:"nama_toko"`
	UrlFoto  string `json:"url_foto"`
	UserID   uint   `json:"user_id,omitempty"` // Hanya diisi untuk GetMyToko
}

type PaginatedTokoResponse struct {
	Page  int            `json:"page"`
	Limit int            `json:"limit"`
	Data  []TokoResponse `json:"data"`
}