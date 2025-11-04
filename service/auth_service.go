package service

import (
	"github.com/Debjth19/go-evermos/helpers"
	"github.com/Debjth19/go-evermos/model"
	"github.com/Debjth19/go-evermos/model/web"
	"github.com/Debjth19/go-evermos/repository"

	"errors"
	"time"

	"gorm.io/gorm"
)

// AuthService adalah interface untuk logika bisnis terkait otentikasi
type AuthService interface {
	Register(request web.AuthRegisterRequest) (model.User, error)
	Login(request web.AuthLoginRequest) (model.User, string, error)
}

// authService adalah implementasi dari AuthService
type authService struct {
	repository repository.AuthRepository
}

// NewAuthService membuat instance baru dari authService
func NewAuthService(repository repository.AuthRepository) AuthService {
	return &authService{repository}
}

// Register menangani logika pendaftaran user baru
func (s *authService) Register(request web.AuthRegisterRequest) (model.User, error) {
	// 1. Validasi duplikat email
	emailExists, err := s.repository.CheckEmailExists(request.Email)
	if err != nil {
		return model.User{}, err
	}
	if emailExists {
		return model.User{}, errors.New("Email sudah terdaftar")
	}

	// 2. Validasi duplikat NoTelp
	noTelpExists, err := s.repository.CheckNoTelpExists(request.NoTelp)
	if err != nil {
		return model.User{}, err
	}
	if noTelpExists {
		return model.User{}, errors.New("Nomor telepon sudah terdaftar")
	}

	// 3. Hash password
	hashedPassword, err := helpers.HashPassword(request.KataSandi)
	if err != nil {
		return model.User{}, err
	}

	// 4. Konversi TanggalLahir 
	tanggalLahir, err := time.Parse("02/01/2006", request.TanggalLahir)
	if err != nil {
		return model.User{}, errors.New("Format tanggal lahir tidak valid, gunakan dd/mm/yyyy")
	}

	// 5. Buat struct model.User untuk disimpan
	newUser := model.User{
		Nama:         request.Nama,
		KataSandi:    hashedPassword,
		NoTelp:       request.NoTelp,
		TanggalLahir: tanggalLahir,
		Pekerjaan:    request.Pekerjaan,
		Email:        request.Email,
		IDProvinsi:   request.IDProvinsi,
		IDKota:       request.IDKota,
		Role:         "user", // Role default
	}

	// 6. Panggil repository untuk menyimpan user (dan toko)
	createdUser, err := s.repository.Register(newUser)
	if err != nil {
		return model.User{}, err
	}

	return createdUser, nil
}

// Login menangani logika login user
func (s *authService) Login(request web.AuthLoginRequest) (model.User, string, error) {
	// 1. Cari user berdasarkan NoTelp
	user, err := s.repository.Login(request.NoTelp)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.User{}, "", errors.New("Nomor telepon atau kata sandi salah")
		}
		return model.User{}, "", err
	}

	// 2. Verifikasi password
	passwordMatch := helpers.CheckPasswordHash(request.KataSandi, user.KataSandi)
	if !passwordMatch {
		return model.User{}, "", errors.New("Nomor telepon atau kata sandi salah")
	}

	// 3. Buat JWT Token
	token, err := helpers.GenerateToken(user.ID, user.Role)
	if err != nil {
		return model.User{}, "", err
	}

	// 4. Kembalikan data user dan token
	return user, token, nil
}