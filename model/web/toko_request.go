package web

type TokoUpdateRequest struct {
	NamaToko string `json:"nama_toko" validate:"required"`
}