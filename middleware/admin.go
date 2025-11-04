package middleware

import (
	"github.com/Debjth19/go-evermos/model/web"
	"github.com/gofiber/fiber/v2"
)

// Middleware untuk memeriksa apakah user adalah Admin
func AdminMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Ambil role yang sudah disimpan oleh AuthMiddleware
		role, ok := c.Locals("role").(string)

		if !ok || role != "admin" {
			// Jika bukan admin, tolak akses
			return c.Status(fiber.StatusForbidden).JSON(web.WebResponse{
				Status:  false,
				Message: "Gagal",
				Errors:  "Akses ditolak: Hanya untuk admin",
			})
		}

		// Jika admin, lanjutkan
		return c.Next()
	}
}