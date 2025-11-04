package web

type AuthRegisterRequest struct {
	Nama         string `json:"nama" validate:"required"`
	KataSandi    string `json:"kata_sandi" validate:"required,min=6"`
	NoTelp       string `json:"no_telp" validate:"required"`
	TanggalLahir string `json:"tanggal_Lahir" validate:"required"` // Format "dd/mm/yyyy"
	Pekerjaan    string `json:"pekerjaan" validate:"required"`
	Email        string `json:"email" validate:"required,email"`
	IDProvinsi   string `json:"id_provinsi" validate:"required"`
	IDKota       string `json:"id_kota" validate:"required"`
}

type AuthLoginRequest struct {
	NoTelp    string `json:"no_telp" validate:"required"`
	KataSandi string `json:"kata_sandi" validate:"required"`
}