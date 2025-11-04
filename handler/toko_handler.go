package handler

import (
	"github.com/Debjth19/go-evermos/helpers"
	"github.com/Debjth19/go-evermos/model/web"
	"github.com/Debjth19/go-evermos/service"

	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type TokoHandler interface {
	GetMyToko(c *fiber.Ctx) error
	GetTokoByID(c *fiber.Ctx) error
	GetAllToko(c *fiber.Ctx) error
	UpdateToko(c *fiber.Ctx) error
}

type tokoHandler struct {
	tokoService service.TokoService
}

func NewTokoHandler(tokoService service.TokoService) TokoHandler {
	return &tokoHandler{tokoService: tokoService}
}

// GetMyToko menangani GET /toko/my 
func (h *tokoHandler) GetMyToko(c *fiber.Ctx) error {
	// 1. Ambil user_id dari middleware
	userID := c.Locals("user_id").(uint)

	// 2. Panggil service
	toko, err := h.tokoService.GetMyToko(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(web.WebResponse{
			Status:  false,
			Message: "Gagal",
			Errors:  err.Error(),
		})
	}

	// 3. Buat respons
	response := web.TokoResponse{
		ID:       toko.ID,
		NamaToko: toko.NamaToko,
		UrlFoto:  toko.UrlFoto,
		UserID:   toko.UserID,
	}

	return c.Status(fiber.StatusOK).JSON(web.WebResponse{
		Status:  true,
		Message: "Succeed to GET data",
		Data:    response,
	})
}

// GetTokoByID menangani GET /toko/:id_toko
func (h *tokoHandler) GetTokoByID(c *fiber.Ctx) error {
	// 1. Ambil :id_toko dari URL
	tokoID, err := strconv.Atoi(c.Params("id_toko"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.WebResponse{
			Status:  false,
			Message: "Bad Request",
			Errors:  "ID toko tidak valid",
		})
	}

	// 2. Panggil service
	toko, err := h.tokoService.GetTokoByID(uint(tokoID))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(web.WebResponse{
			Status:  false,
			Message: "Gagal",
			Errors:  err.Error(),
		})
	}

	// 3. Buat respons
	response := web.TokoResponse{
		ID:       toko.ID,
		NamaToko: toko.NamaToko,
		UrlFoto:  toko.UrlFoto,
	}

	return c.Status(fiber.StatusOK).JSON(web.WebResponse{
		Status:  true,
		Message: "Succeed to GET data",
		Data:    response,
	})
}

// GetAllToko menangani GET /toko 
func (h *tokoHandler) GetAllToko(c *fiber.Ctx) error {
	pagination := helpers.GeneratePagination(c)

	// 2. Ambil query param untuk filter nama 
	search := c.Query("nama") 
	tokos, err := h.tokoService.GetAllToko(pagination, search)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(web.WebResponse{
			Status:  false,
			Message: "Server Error",
			Errors:  err.Error(),
		})
	}

	// 4. Buat respons
	var response []web.TokoResponse
	for _, toko := range tokos {
		response = append(response, web.TokoResponse{
			ID:       toko.ID,
			NamaToko: toko.NamaToko,
			UrlFoto:  toko.UrlFoto,
		})
	}

	// 5. Buat respons paginasi
	paginatedResponse := web.PaginatedTokoResponse{
		Page:  pagination.Page,
		Limit: pagination.Limit,
		Data:  response,
	}

	return c.Status(fiber.StatusOK).JSON(web.WebResponse{
		Status:  true,
		Message: "Succeed to GET data",
		Data:    paginatedResponse,
	})
}

// UpdateToko menangani PUT /toko/:id_toko 
func (h *tokoHandler) UpdateToko(c *fiber.Ctx) error {
	// 1. Ambil user_id dari middleware
	userID := c.Locals("user_id").(uint)

	// 2. Ambil :id_toko dari URL
	tokoID, err := strconv.Atoi(c.Params("id_toko"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.WebResponse{
			Status:  false,
			Message: "Bad Request",
			Errors:  "ID toko tidak valid",
		})
	}

	// 3. Ambil data form-data
	request := web.TokoUpdateRequest{
		NamaToko: c.FormValue("nama_toko"),
	}

	// 4. Ambil file dari form-data 
	file, err := c.FormFile("photo") 
	if err != nil {
		if err.Error() != "there is no uploaded file associated with the given key" {
			return c.Status(fiber.StatusInternalServerError).JSON(web.WebResponse{
				Status:  false,
				Message: "Server Error",
				Errors:  "Gagal memproses file: " + err.Error(),
			})
		}
		file = nil
	}
	
	// 5. Panggil service
	_, err = h.tokoService.UpdateToko(userID, uint(tokoID), request, file)
	if err != nil {
		if strings.Contains(err.Error(), "Akses ditolak") {
			return c.Status(fiber.StatusForbidden).JSON(web.WebResponse{ // 403 Forbidden
				Status:  false,
				Message: "Gagal",
				Errors:  err.Error(),
			})
		}
		if strings.Contains(err.Error(), "tidak ditemukan") {
			return c.Status(fiber.StatusNotFound).JSON(web.WebResponse{ // 404 Not Found
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

	// 6. Kirim respons sukses
	return c.Status(fiber.StatusOK).JSON(web.WebResponse{
		Status:  true,
		Message: "Succeed to UPDATE data",
		Data:    "Update toko succeed",
	})
}