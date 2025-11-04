package service

import (
	"github.com/Debjth19/go-evermos/helpers"
	"github.com/Debjth19/go-evermos/model"
	"github.com/Debjth19/go-evermos/model/web"
	"github.com/Debjth19/go-evermos/repository"

	"errors"
	"time"
)

type UserService interface {
	GetProfile(userID uint) (model.User, error)
	UpdateProfile(userID uint, request web.UserUpdateRequest) (model.User, error)
}

type userService struct {
	userRepository repository.UserRepository
	authRepository repository.AuthRepository 
}

func NewUserService(userRepo repository.UserRepository, authRepo repository.AuthRepository) UserService {
	return &userService{
		userRepository: userRepo,
		authRepository: authRepo,
	}
}

func (s *userService) GetProfile(userID uint) (model.User, error) {
	user, err := s.userRepository.FindByID(userID)
	if err != nil {
		return user, errors.New("User tidak ditemukan")
	}
	return user, nil
}

func (s *userService) UpdateProfile(userID uint, request web.UserUpdateRequest) (model.User, error) {
	// 1. Dapatkan data user saat ini
	user, err := s.userRepository.FindByID(userID)
	if err != nil {
		return user, errors.New("User tidak ditemukan")
	}

	// 2. Cek duplikat Email (jika email diubah)
	if request.Email != "" && request.Email != user.Email {
		emailExists, err := s.authRepository.CheckEmailExists(request.Email)
		if err != nil {
			return user, err
		}
		if emailExists {
			return user, errors.New("Email sudah terdaftar")
		}
		user.Email = request.Email
	}

	// 3. Cek duplikat NoTelp (jika no_telp diubah)
	if request.NoTelp != "" && request.NoTelp != user.NoTelp {
		noTelpExists, err := s.authRepository.CheckNoTelpExists(request.NoTelp)
		if err != nil {
			return user, err
		}
		if noTelpExists {
			return user, errors.New("Nomor telepon sudah terdaftar")
		}
		user.NoTelp = request.NoTelp
	}

	// 4. Update field lainnya
	if request.Nama != "" {
		user.Nama = request.Nama
	}
	if request.Pekerjaan != "" {
		user.Pekerjaan = request.Pekerjaan
	}
	if request.IDProvinsi != "" {
		user.IDProvinsi = request.IDProvinsi
	}
	if request.IDKota != "" {
		user.IDKota = request.IDKota
	}
	
	// 5. Update Tanggal Lahir (jika diisi)
	if request.TanggalLahir != "" {
		tanggalLahir, err := time.Parse("02/01/2006", request.TanggalLahir)
		if err != nil {
			return user, errors.New("Format tanggal lahir tidak valid, gunakan dd/mm/yyyy")
		}
		user.TanggalLahir = tanggalLahir
	}

	// 6. Update Password (jika diisi)
	if request.KataSandi != "" {
		hashedPassword, err := helpers.HashPassword(request.KataSandi)
		if err != nil {
			return user, err
		}
		user.KataSandi = hashedPassword
	}

	// 7. Simpan perubahan ke DB
	updatedUser, err := s.userRepository.Update(user)
	if err != nil {
		return updatedUser, err
	}

	return updatedUser, nil
}