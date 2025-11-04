package handler

import (
	"github.com/Debjth19/go-evermos/model/web"
	"github.com/Debjth19/go-evermos/service"

	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type AlamatHandler interface {
	CreateAlamat(c *fiber.Ctx) error
	GetAllAlamat(c *fiber.Ctx) error
	GetAlamatByID(c *fiber.Ctx) error
	UpdateAlamat(c *fiber.Ctx) error
	DeleteAlamat(c *fiber.Ctx) error
}

type alamatHandler struct {
	alamatService service.AlamatService
}

func NewAlamatHandler(alamatService service.AlamatService) AlamatHandler {
	return &alamatHandler{alamatService}
}

// CreateAlamat menangani POST /user/alamat
func (h *alamatHandler) CreateAlamat(c *fiber.Ctx) error {
	// 1. Ambil user_id dari middleware
	userID := c.Locals("user_id").(uint)

	// 2. Parse request body
	var request web.AlamatCreateRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.WebResponse{
			Status:  false,
			Message: "Bad Request",
			Errors:  err.Error(),
		})
	}

	// 3. Panggil service
	newAlamat, err := h.alamatService.CreateAlamat(userID, request)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(web.WebResponse{
			Status:  false,
			Message: "Server Error",
			Errors:  err.Error(),
		})
	}

	// 4. Buat response
	response := web.AlamatResponse{
		ID:           newAlamat.ID,
		JudulAlamat:  newAlamat.JudulAlamat,
		NamaPenerima: newAlamat.NamaPenerima,
		NoTelp:       newAlamat.NoTelp,
		DetailAlamat: newAlamat.DetailAlamat,
		UserID:       newAlamat.UserID,
	}

	return c.Status(fiber.StatusOK).JSON(web.WebResponse{
		Status:  true,
		Message: "Succeed to POST data",
		Data:    response,
	})
}

// GetAllAlamat menangani GET /user/alamat
func (h *alamatHandler) GetAllAlamat(c *fiber.Ctx) error {
	// 1. Ambil user_id dari middleware
	userID := c.Locals("user_id").(uint)

	// 2. Panggil service
	alamats, err := h.alamatService.GetAllAlamat(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(web.WebResponse{
			Status:  false,
			Message: "Server Error",
			Errors:  err.Error(),
		})
	}

	// 3. Buat response
	var response []web.AlamatResponse
	for _, alamat := range alamats {
		response = append(response, web.AlamatResponse{
			ID:           alamat.ID,
			JudulAlamat:  alamat.JudulAlamat,
			NamaPenerima: alamat.NamaPenerima,
			NoTelp:       alamat.NoTelp,
			DetailAlamat: alamat.DetailAlamat,
			UserID:       alamat.UserID,
		})
	}

	return c.Status(fiber.StatusOK).JSON(web.WebResponse{
		Status:  true,
		Message: "Succeed to GET data",
		Data:    response,
	})
}

// GetAlamatByID menangani GET /user/alamat/:id
func (h *alamatHandler) GetAlamatByID(c *fiber.Ctx) error {
	// 1. Ambil user_id dari middleware
	userID := c.Locals("user_id").(uint)

	// 2. Ambil :id dari URL parameter
	alamatID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.WebResponse{
			Status:  false,
			Message: "Bad Request",
			Errors:  "ID alamat tidak valid",
		})
	}

	// 3. Panggil service
	alamat, err := h.alamatService.GetAlamatByID(userID, uint(alamatID))
	if err != nil {
		// Cek jika errornya adalah "Akses ditolak"
		if strings.Contains(err.Error(), "Akses ditolak") {
			return c.Status(fiber.StatusForbidden).JSON(web.WebResponse{ // 403 Forbidden
				Status:  false,
				Message: "Gagal",
				Errors:  err.Error(),
			})
		}
		// Cek jika errornya adalah "Tidak ditemukan"
		if strings.Contains(err.Error(), "tidak ditemukan") {
			return c.Status(fiber.StatusNotFound).JSON(web.WebResponse{ // 404 Not Found
				Status:  false,
				Message: "Gagal",
				Errors:  err.Error(),
			})
		}
		// Error lainnya
		return c.Status(fiber.StatusInternalServerError).JSON(web.WebResponse{
			Status:  false,
			Message: "Server Error",
			Errors:  err.Error(),
		})
	}

	// 4. Buat response
	response := web.AlamatResponse{
		ID:           alamat.ID,
		JudulAlamat:  alamat.JudulAlamat,
		NamaPenerima: alamat.NamaPenerima,
		NoTelp:       alamat.NoTelp,
		DetailAlamat: alamat.DetailAlamat,
		UserID:       alamat.UserID,
	}

	return c.Status(fiber.StatusOK).JSON(web.WebResponse{
		Status:  true,
		Message: "Succeed to GET data",
		Data:    response,
	})
}

// UpdateAlamat menangani PUT /user/alamat/:id
func (h *alamatHandler) UpdateAlamat(c *fiber.Ctx) error {
	// 1. Ambil user_id dari middleware
	userID := c.Locals("user_id").(uint)

	// 2. Ambil :id dari URL parameter
	alamatID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.WebResponse{
			Status:  false,
			Message: "Bad Request",
			Errors:  "ID alamat tidak valid",
		})
	}

	// 3. Parse request body
	var request web.AlamatUpdateRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.WebResponse{
			Status:  false,
			Message: "Bad Request",
			Errors:  err.Error(),
		})
	}

	// 4. Panggil service
	_, err = h.alamatService.UpdateAlamat(userID, uint(alamatID), request)
	if err != nil {
		// Cek error seperti di GetAlamatByID
		if strings.Contains(err.Error(), "Akses ditolak") {
			return c.Status(fiber.StatusForbidden).JSON(web.WebResponse{
				Status:  false,
				Message: "Gagal",
				Errors:  err.Error(),
			})
		}
		if strings.Contains(err.Error(), "tidak ditemukan") {
			return c.Status(fiber.StatusNotFound).JSON(web.WebResponse{
				Status:  false,
				Message: "Gagal",
				Errors:  err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(web.WebResponse{
			Status:  false,
			Message: "Server Error",
			Errors:  err.Error(),
		})
	}

	// 5. Kirim respons sukses
	return c.Status(fiber.StatusOK).JSON(web.WebResponse{
		Status:  true,
		Message: "Succeed to UPDATE data",
		Data:    "Update alamat succeed",
	})
}

// DeleteAlamat menangani DELETE /user/alamat/:id
func (h *alamatHandler) DeleteAlamat(c *fiber.Ctx) error {
	// 1. Ambil user_id dari middleware
	userID := c.Locals("user_id").(uint)

	// 2. Ambil :id dari URL parameter
	alamatID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.WebResponse{
			Status:  false,
			Message: "Bad Request",
			Errors:  "ID alamat tidak valid",
		})
	}

	// 3. Panggil service
	err = h.alamatService.DeleteAlamat(userID, uint(alamatID))
	if err != nil {
		if strings.Contains(err.Error(), "Akses ditolak") {
			return c.Status(fiber.StatusForbidden).JSON(web.WebResponse{
				Status:  false,
				Message: "Gagal",
				Errors:  err.Error(),
			})
		}
		if strings.Contains(err.Error(), "tidak ditemukan") {
			return c.Status(fiber.StatusNotFound).JSON(web.WebResponse{
				Status:  false,
				Message: "Gagal",
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
		Message: "Succeed to DELETE data",
		Data:    "Delete alamat succeed",
	})
}