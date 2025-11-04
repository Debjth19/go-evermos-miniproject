package handler

import (
	"github.com/Debjth19/go-evermos/model/web"
	"github.com/Debjth19/go-evermos/service"

	"strings"

	"github.com/gofiber/fiber/v2"
)

type UserHandler interface {
	GetProfile(c *fiber.Ctx) error
	UpdateProfile(c *fiber.Ctx) error
}

type userHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) UserHandler {
	return &userHandler{userService}
}

// GetProfile menangani GET /user
func (h *userHandler) GetProfile(c *fiber.Ctx) error {
	// 1. Ambil user_id dari c.Locals 
	userID := c.Locals("user_id").(uint)

	// 2. Panggil service
	user, err := h.userService.GetProfile(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(web.WebResponse{
			Status:  false,
			Message: "Gagal",
			Errors:  err.Error(),
		})
	}

	// 3. Buat respons
	response := web.UserResponse{
		ID:           user.ID,
		Nama:         user.Nama,
		NoTelp:       user.NoTelp,
		TanggalLahir: user.TanggalLahir.Format("02/01/2006"),
		Pekerjaan:    user.Pekerjaan,
		Email:        user.Email,
		IDProvinsi:   user.IDProvinsi,
		IDKota:       user.IDKota,
		Role:         user.Role,
	}

	return c.Status(fiber.StatusOK).JSON(web.WebResponse{
		Status:  true,
		Message: "Succeed to GET data",
		Data:    response,
	})
}

// UpdateProfile menangani PUT /user
func (h *userHandler) UpdateProfile(c *fiber.Ctx) error {
	// 1. Ambil user_id dari c.Locals
	userID := c.Locals("user_id").(uint)

	// 2. Parse request body
	var request web.UserUpdateRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.WebResponse{
			Status:  false,
			Message: "Bad Request",
			Errors:  err.Error(),
		})
	}

	// 3. Panggil service
	_, err := h.userService.UpdateProfile(userID, request)
	if err != nil {
		if strings.Contains(err.Error(), "terdaftar") || strings.Contains(err.Error(), "format tanggal") {
			return c.Status(fiber.StatusConflict).JSON(web.WebResponse{ // 409 Conflict
				Status:  false,
				Message: "Update Gagal",
				Errors:  err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(web.WebResponse{
			Status:  false,
			Message: "Server Error",
			Errors:  err.Error(),
		})
	}

	// 4. Kirim respons sukses
	return c.Status(fiber.StatusOK).JSON(web.WebResponse{
		Status:  true,
		Message: "Succeed to UPDATE data",
		Data:    "Update profile succeed",
	})
}