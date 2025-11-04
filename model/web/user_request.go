package web

type UserUpdateRequest struct {
	Nama         string `json:"nama"`
	KataSandi    string `json:"kata_sandi"` 
	NoTelp       string `json:"no_telp"`
	TanggalLahir string `json:"tanggal_Lahir"` 
	Pekerjaan    string `json:"pekerjaan"`
	Email        string `json:"email"`
	IDProvinsi   string `json:"id_provinsi"`
	IDKota       string `json:"id_kota"`
}