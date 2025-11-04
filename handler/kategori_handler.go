package handler

import (
	"github.com/Debjth19/go-evermos/model/web"
	"github.com/Debjth19/go-evermos/service"

	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type KategoriHandler interface {
	CreateKategori(c *fiber.Ctx) error
	GetAllKategori(c *fiber.Ctx) error
	GetKategoriByID(c *fiber.Ctx) error
	UpdateKategori(c *fiber.Ctx) error
	DeleteKategori(c *fiber.Ctx) error
}

type kategoriHandler struct {
	kategoriService service.KategoriService
}

func NewKategoriHandler(kategoriService service.KategoriService) KategoriHandler {
	return &kategoriHandler{kategoriService}
}

// CreateKategori menangani POST /category
func (h *kategoriHandler) CreateKategori(c *fiber.Ctx) error {
	var request web.KategoriCreateRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.WebResponse{
			Status:  false,
			Message: "Bad Request",
			Errors:  err.Error(),
		})
	}

	kategori, err := h.kategoriService.CreateKategori(request)
	if err != nil {
		if strings.Contains(err.Error(), "sudah ada") {
			return c.Status(fiber.StatusConflict).JSON(web.WebResponse{ // 409 Conflict
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

	return c.Status(fiber.StatusOK).JSON(web.WebResponse{
		Status:  true,
		Message: "Succeed to POST data",
		Data:    kategori.ID,
	})
}

// GetAllKategori menangani GET /category
func (h *kategoriHandler) GetAllKategori(c *fiber.Ctx) error {
	kategoris, err := h.kategoriService.GetAllKategori()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(web.WebResponse{
			Status:  false,
			Message: "Server Error",
			Errors:  err.Error(),
		})
	}

	var response []web.KategoriResponse
	for _, k := range kategoris {
		response = append(response, web.KategoriResponse{
			ID:           k.ID,
			NamaCategory: k.NamaCategory,
		})
	}

	return c.Status(fiber.StatusOK).JSON(web.WebResponse{
		Status:  true,
		Message: "Succeed to GET data",
		Data:    response,
	})
}

// GetKategoriByID menangani GET /category/:id
func (h *kategoriHandler) GetKategoriByID(c *fiber.Ctx) error {
	kategoriID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.WebResponse{
			Status:  false,
			Message: "Bad Request",
			Errors:  "ID kategori tidak valid",
		})
	}

	kategori, err := h.kategoriService.GetKategoriByID(uint(kategoriID))
	if err != nil {
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

	response := web.KategoriResponse{
		ID:           kategori.ID,
		NamaCategory: kategori.NamaCategory,
	}

	return c.Status(fiber.StatusOK).JSON(web.WebResponse{
		Status:  true,
		Message: "Succeed to GET data",
		Data:    response,
	})
}

// UpdateKategori menangani PUT /category/:id
func (h *kategoriHandler) UpdateKategori(c *fiber.Ctx) error {
	kategoriID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.WebResponse{
			Status:  false,
			Message: "Bad Request",
			Errors:  "ID kategori tidak valid",
		})
	}

	var request web.KategoriUpdateRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.WebResponse{
			Status:  false,
			Message: "Bad Request",
			Errors:  err.Error(),
		})
	}

	_, err = h.kategoriService.UpdateKategori(uint(kategoriID), request)
	if err != nil {
		if strings.Contains(err.Error(), "tidak ditemukan") {
			return c.Status(fiber.StatusNotFound).JSON(web.WebResponse{
				Status:  false,
				Message: "Gagal",
				Errors:  err.Error(),
			})
		}
		if strings.Contains(err.Error(), "sudah ada") {
			return c.Status(fiber.StatusConflict).JSON(web.WebResponse{
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

	return c.Status(fiber.StatusOK).JSON(web.WebResponse{
		Status:  true,
		Message: "Succeed to UPDATE data",
		Data:    "",
	})
}

// DeleteKategori menangani DELETE /category/:id
func (h *kategoriHandler) DeleteKategori(c *fiber.Ctx) error {
	kategoriID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.WebResponse{
			Status:  false,
			Message: "Bad Request",
			Errors:  "ID kategori tidak valid",
		})
	}

	err = h.kategoriService.DeleteKategori(uint(kategoriID))
	if err != nil {
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

	return c.Status(fiber.StatusOK).JSON(web.WebResponse{
		Status:  true,
		Message: "Succeed to DELETE data",
		Data:    "",
	})
}