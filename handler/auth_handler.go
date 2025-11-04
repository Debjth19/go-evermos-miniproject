package handler

import (
	"github.com/Debjth19/go-evermos/model/web"
	"github.com/Debjth19/go-evermos/service"

	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler interface {
	Register(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
}

type authHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) AuthHandler {
	return &authHandler{authService}
}

// Register menangani permintaan POST /auth/register
func (h *authHandler) Register(c *fiber.Ctx) error {
	// 1. Parse request body ke struct AuthRegisterRequest
	var request web.AuthRegisterRequest
	if err := c.BodyParser(&request); err != nil {
		// Jika parsing gagal 
		return c.Status(fiber.StatusBadRequest).JSON(web.WebResponse{
			Status:  false,
			Message: "Bad Request: " + err.Error(),
			Errors:  err.Error(),
		})
	}

	// 2. Panggil service Register
	user, err := h.authService.Register(request)
	if err != nil {
		// Jika ada error dari service (misal: email duplikat)
		if strings.Contains(err.Error(), "sudah terdaftar") || strings.Contains(err.Error(), "format tanggal") {
			return c.Status(fiber.StatusConflict).JSON(web.WebResponse{ // 409 Conflict
				Status:  false,
				Message: "Registrasi Gagal",
				Errors:  err.Error(),
			})
		}

		// Untuk error lainnya (misal: error database)
		return c.Status(fiber.StatusInternalServerError).JSON(web.WebResponse{ // 500 Internal Server Error
			Status:  false,
			Message: "Server Error",
			Errors:  err.Error(),
		})
	}

	// 3. Kirim respons sukses
	return c.Status(fiber.StatusOK).JSON(web.WebResponse{
		Status:  true,
		Message: "Succeed to POST data",
		Data:    fmt.Sprintf("Register Succeed for user: %s (ID: %d)", user.Nama, user.ID), // Kita buat lebih jelas
	})
}

// Login menangani permintaan POST /auth/login
func (h *authHandler) Login(c *fiber.Ctx) error {
	// 1. Parse request body ke struct AuthLoginRequest
	var request web.AuthLoginRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.WebResponse{
			Status:  false,
			Message: "Bad Request: " + err.Error(),
			Errors:  err.Error(),
		})
	}

	// 2. Panggil service Login
	user, token, err := h.authService.Login(request)
	if err != nil {
		// Jika errornya "No Telp atau kata sandi salah"
		if strings.Contains(err.Error(), "salah") {
			return c.Status(fiber.StatusUnauthorized).JSON(web.WebResponse{ // 401 Unauthorized
				Status:  false,
				Message: "Failed to POST data",
				Errors:  "Nomor telepon atau kata sandi salah",
			})
		}
		
		// Error internal lainnya
		return c.Status(fiber.StatusInternalServerError).JSON(web.WebResponse{
			Status:  false,
			Message: "Server Error",
			Errors:  err.Error(),
		})
	}

	// 3. Buat respons body sesuai Postman
	loginResponse := web.AuthLoginResponse{
		Nama:         user.Nama,
		NoTelp:       user.NoTelp,
		TanggalLahir: user.TanggalLahir.Format("02/01/2006"), 
		Tentang:      "", 
		Pekerjaan:    user.Pekerjaan,
		Email:        user.Email,
		IDProvinsi:   user.IDProvinsi, 
		IDKota:       user.IDKota,    
		Token:        token,
	}

	// 4. Kirim respons sukses
	return c.Status(fiber.StatusOK).JSON(web.WebResponse{
		Status:  true,
		Message: "Succeed to POST data",
		Data:    loginResponse,
	})
}