package helpers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type Pagination struct {
	Limit int `json:"limit"`
	Page  int `json:"page"`
	
}

// GeneratePagination mem-parsing query params dari URL
func GeneratePagination(c *fiber.Ctx) Pagination {
	// Set nilai default
	limit := 10
	page := 1

	// Ambil query dan konversi ke integer
	if l, err := strconv.Atoi(c.Query("limit")); err == nil && l > 0 {
		limit = l
	}

	if p, err := strconv.Atoi(c.Query("page")); err == nil && p > 0 {
		page = p
	}

	return Pagination{
		Limit: limit,
		Page:  page,
	}
}