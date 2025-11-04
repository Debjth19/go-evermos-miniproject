package handler

import (
	"github.com/Debjth19/go-evermos/model"
	"github.com/Debjth19/go-evermos/model/web"
	"github.com/Debjth19/go-evermos/service"

	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type TransaksiHandler interface {
	CreateTransaksi(c *fiber.Ctx) error
	GetMyTransactions(c *fiber.Ctx) error
	GetMyTransactionByID(c *fiber.Ctx) error
}

type transaksiHandler struct {
	transaksiService service.TransaksiService
}

func NewTransaksiHandler(transaksiService service.TransaksiService) TransaksiHandler {
	return &transaksiHandler{transaksiService: transaksiService}
}

// CreateTransaksi menangani POST /trx
func (h *transaksiHandler) CreateTransaksi(c *fiber.Ctx) error {
	// 1. Ambil user_id dari middleware
	userID := c.Locals("user_id").(uint)

	// 2. Parse request body
	var request web.TransaksiCreateRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.WebResponse{
			Status:  false,
			Message: "Bad Request",
			Errors:  err.Error(),
		})
	}

	// 3. Panggil service
	transaksi, err := h.transaksiService.CreateTransaksi(userID, request)
	if err != nil {
		if strings.Contains(err.Error(), "Stok tidak mencukupi") ||
			strings.Contains(err.Error(), "tidak ditemukan") ||
			strings.Contains(err.Error(), "Akses ditolak") {
			return c.Status(fiber.StatusBadRequest).JSON(web.WebResponse{ // 400 Bad Request
				Status:  false,
				Message: "Gagal membuat transaksi",
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
		Data:    transaksi.ID,
	})
}

// GetMyTransactions menangani GET /trx
func (h *transaksiHandler) GetMyTransactions(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	transaksis, err := h.transaksiService.GetMyTransactions(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(web.WebResponse{
			Status:  false,
			Message: "Server Error",
			Errors:  err.Error(),
		})
	}

	// Buat respons
	var response []web.TransaksiResponse
	for _, trx := range transaksis {
		response = append(response, mapTransaksiToResponse(trx))
	}

	return c.Status(fiber.StatusOK).JSON(web.WebResponse{
		Status:  true,
		Message: "Succeed to GET data",
		Data:    web.PaginatedTransaksiResponse{Data: response}, // Sesuai Postman, ada 'data' di dalamnya
	})
}

// GetMyTransactionByID menangani GET /trx/:id
func (h *transaksiHandler) GetMyTransactionByID(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	trxID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.WebResponse{
			Status:  false,
			Message: "Bad Request",
			Errors:  "ID transaksi tidak valid",
		})
	}

	transaksi, err := h.transaksiService.GetMyTransactionByID(userID, uint(trxID))
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
		Data:    mapTransaksiToResponse(transaksi),
	})
}

// --- Helper Mapping ---

func mapTransaksiToResponse(t model.Transaksi) web.TransaksiResponse {
	return web.TransaksiResponse{
		ID:          t.ID,
		HargaTotal:  t.HargaTotal,
		KodeInvoice: t.KodeInvoice,
		MethodBayar: t.MethodBayar,
		AlamatKirim: mapAlamatToResponse(t.Alamat),
		DetailTrx:   mapDetailTrxToResponse(t.DetailTransaksi),
	}
}

func mapAlamatToResponse(a model.Alamat) web.AlamatResponse {
	return web.AlamatResponse{
		ID:           a.ID,
		JudulAlamat:  a.JudulAlamat,
		NamaPenerima: a.NamaPenerima,
		NoTelp:       a.NoTelp,
		DetailAlamat: a.DetailAlamat,
		UserID:       a.UserID,
	}
}

func mapDetailTrxToResponse(details []model.DetailTransaksi) []web.DetailTransaksiResponse {
	var response []web.DetailTransaksiResponse
	for _, d := range details {
		response = append(response, web.DetailTransaksiResponse{
			Produk:     MapProdukToResponse(d.Produk),
			Toko:       MapTokoToResponse(d.Toko),
			Kuantitas:  d.Kuantitas,
			HargaTotal: d.HargaTotal,
		})
	}
	return response
}
