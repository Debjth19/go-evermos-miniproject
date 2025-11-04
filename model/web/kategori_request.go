package web

type KategoriCreateRequest struct {
	NamaCategory string `json:"nama_category" validate:"required"`
}

type KategoriUpdateRequest struct {
	NamaCategory string `json:"nama_category" validate:"required"`
}