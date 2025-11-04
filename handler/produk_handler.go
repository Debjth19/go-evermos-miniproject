package handler

import (
	"errors"

	"github.com/Debjth19/go-evermos/helpers"
	"github.com/Debjth19/go-evermos/model"
	"github.com/Debjth19/go-evermos/model/web"
	"github.com/Debjth19/go-evermos/service"

	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type ProdukHandler interface {
	CreateProduk(c *fiber.Ctx) error
	GetAllProduk(c *fiber.Ctx) error
	GetProdukByID(c *fiber.Ctx) error
	UpdateProduk(c *fiber.Ctx) error
	DeleteProduk(c *fiber.Ctx) error
}

type produkHandler struct {
	produkService service.ProdukService
}

func NewProdukHandler(produkService service.ProdukService) ProdukHandler {
	return &produkHandler{produkService: produkService}
}

func parseProdukCreateRequest(c *fiber.Ctx) (web.ProdukCreateRequest, error) {
	request := web.ProdukCreateRequest{
		NamaProduk: c.FormValue("nama_produk"),
		Deskripsi:  c.FormValue("deskripsi"),
	}

	catID, err := strconv.Atoi(c.FormValue("category_id"))
	if err != nil {
		return request, errors.New("category_id tidak valid")
	}
	request.CategoryID = uint(catID)

	hargaReseler, err := strconv.Atoi(c.FormValue("harga_reseller"))
	if err != nil {
		return request, errors.New("harga_reseller tidak valid")
	}
	request.HargaReseler = uint(hargaReseler)

	hargaKonsumen, err := strconv.Atoi(c.FormValue("harga_konsumen"))
	if err != nil {
		return request, errors.New("harga_konsumen tidak valid")
	}
	request.HargaKonsumen = uint(hargaKonsumen)

	stok, err := strconv.Atoi(c.FormValue("stok"))
	if err != nil {
		return request, errors.New("stok tidak valid")
	}
	request.Stok = uint(stok)

	return request, nil
}

// CreateProduk menangani POST /product
func (h *produkHandler) CreateProduk(c *fiber.Ctx) error {
	// 1. Ambil user_id dari middleware
	userID := c.Locals("user_id").(uint)

	// 2. Parse request form-data
	request, err := parseProdukCreateRequest(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.WebResponse{
			Status:  false,
			Message: "Bad Request",
			Errors:  err.Error(),
		})
	}

	// 3. Ambil file (photos)
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.WebResponse{
			Status:  false,
			Message: "Bad Request",
			Errors:  "Gagal memproses form: " + err.Error(),
		})
	}
	files := form.File["photos"] 

	// 4. Panggil service
	newProduk, err := h.produkService.CreateProduk(userID, request, files)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(web.WebResponse{
			Status:  false,
			Message: "Server Error",
			Errors:  err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(web.WebResponse{
		Status:  true,
		Message: "Succeed to POST data",
		Data:    newProduk.ID,
	})
}

// GetAllProduk menangani GET /product 
func (h *produkHandler) GetAllProduk(c *fiber.Ctx) error {
	// 1. Ambil query params untuk pagination
	pagination := helpers.GeneratePagination(c)

	// 2. Ambil semua query params untuk filter
	filterParams := map[string]string{
		"nama_produk": c.Query("nama_produk"),
		"category_id": c.Query("category_id"),
		"toko_id":     c.Query("toko_id"),
		"min_harga":   c.Query("min_harga"),
		"max_harga":   c.Query("max_harga"),
	}

	// 3. Panggil service
	produks, err := h.produkService.GetAllProduk(pagination, filterParams)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(web.WebResponse{
			Status:  false,
			Message: "Server Error",
			Errors:  err.Error(),
		})
	}

	// 4. Buat respons
	var response []web.ProdukResponse
	for _, p := range produks {
		response = append(response, MapProdukToResponse(p))
	}

	paginatedResponse := web.PaginatedProdukResponse{
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

// GetProdukByID menangani GET /product/:id 
func (h *produkHandler) GetProdukByID(c *fiber.Ctx) error {
	produkID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.WebResponse{
			Status:  false,
			Message: "Bad Request",
			Errors:  "ID produk tidak valid",
		})
	}

	produk, err := h.produkService.GetProdukByID(uint(produkID))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(web.WebResponse{
			Status:  false,
			Message: "Gagal",
			Errors:  err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(web.WebResponse{
		Status:  true,
		Message: "Succeed to GET data",
		Data:    MapProdukToResponse(produk),
	})
}

func parseProdukUpdateRequest(c *fiber.Ctx) web.ProdukUpdateRequest {
	request := web.ProdukUpdateRequest{
		NamaProduk: c.FormValue("nama_produk"),
		Deskripsi:  c.FormValue("deskripsi"),
	}
	if catID, err := strconv.Atoi(c.FormValue("category_id")); err == nil {
		request.CategoryID = uint(catID)
	}
	if hargaReseler, err := strconv.Atoi(c.FormValue("harga_reseller")); err == nil {
		request.HargaReseler = uint(hargaReseler)
	}
	if hargaKonsumen, err := strconv.Atoi(c.FormValue("harga_konsumen")); err == nil {
		request.HargaKonsumen = uint(hargaKonsumen)
	}
	if stok, err := strconv.Atoi(c.FormValue("stok")); err == nil {
		request.Stok = uint(stok)
	}
	return request
}

// UpdateProduk menangani PUT /product/:id
func (h *produkHandler) UpdateProduk(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	produkID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.WebResponse{
			Status:  false,
			Message: "Bad Request",
			Errors:  "ID produk tidak valid",
		})
	}

	request := parseProdukUpdateRequest(c)

	form, _ := c.MultipartForm()
	files := form.File["photos"] 

	_, err = h.produkService.UpdateProduk(userID, uint(produkID), request, files)
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

	return c.Status(fiber.StatusOK).JSON(web.WebResponse{
		Status:  true,
		Message: "Succeed to UPDATE data",
		Data:    "",
	})
}

// DeleteProduk menangani DELETE /product/:id 
func (h *produkHandler) DeleteProduk(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	produkID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.WebResponse{
			Status:  false,
			Message: "Bad Request",
			Errors:  "ID produk tidak valid",
		})
	}

	err = h.produkService.DeleteProduk(userID, uint(produkID))
	if err != nil {
		if strings.Contains(err.Error(), "Akses ditolak") {
			return c.Status(fiber.StatusForbidden).JSON(web.WebResponse{Status: false, Message: "Gagal", Errors: err.Error()})
		}
		if strings.Contains(err.Error(), "tidak ditemukan") {
			return c.Status(fiber.StatusNotFound).JSON(web.WebResponse{Status: false, Message: "Gagal", Errors: err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(web.WebResponse{Status: false, Message: "Server Error", Errors: err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(web.WebResponse{
		Status:  true,
		Message: "Succeed to DELETE data",
		Data:    "",
	})
}

// --- Helper untuk mapping ---
func MapProdukToResponse(p model.Produk) web.ProdukResponse {
	return web.ProdukResponse{
		ID:            p.ID,
		NamaProduk:    p.NamaProduk,
		Slug:          p.Slug,
		HargaReseler:  p.HargaReseler,
		HargaKonsumen: p.HargaKonsumen,
		Stok:          p.Stok,
		Deskripsi:     p.Deskripsi,
		Toko: web.TokoResponse{
			ID:       p.Toko.ID,
			NamaToko: p.Toko.NamaToko,
			UrlFoto:  p.Toko.UrlFoto,
		},
		Category: web.KategoriResponse{
			ID:           p.Category.ID,
			NamaCategory: p.Category.NamaCategory,
		},
		Photos: MapFotosToResponse(p.FotoProduk),
	}
}

func MapTokoToResponse(t model.Toko) web.TokoResponse {
    return web.TokoResponse{
        ID:       t.ID,
        NamaToko: t.NamaToko,
        UrlFoto:  t.UrlFoto,
    }
}

func MapKategoriToResponse(k model.Kategori) web.KategoriResponse {
    return web.KategoriResponse{
        ID:           k.ID,
        NamaCategory: k.NamaCategory,
    }
}

func MapFotosToResponse(fotos []model.FotoProduk) []web.FotoProdukResponse {
	var response []web.FotoProdukResponse
	for _, f := range fotos {
		response = append(response, web.FotoProdukResponse{
			ID:        f.ID,
			ProductID: f.ProductID,
			Url:       f.Url,
		})
	}
	return response
}