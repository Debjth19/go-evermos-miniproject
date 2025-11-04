package web

type AuthLoginResponse struct {
	Nama         string `json:"nama"`
	NoTelp       string `json:"no_telp"`
	TanggalLahir string `json:"tanggal_Lahir"`
	Tentang      string `json:"tentang"` 
	Pekerjaan    string `json:"pekerjaan"`
	Email        string `json:"email"`
	IDProvinsi   string `json:"id_provinsi"`
	IDKota       string `json:"id_kota"`
	Token        string `json:"token"`
}

type WebResponse struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Errors  interface{} `json:"errors,omitempty"` 
	Data    interface{} `json:"data,omitempty"`   
}