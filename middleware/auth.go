package middleware

import (
	"github.com/Debjth19/go-evermos/helpers"
	"github.com/Debjth19/go-evermos/model/web"

	"github.com/gofiber/fiber/v2"
)

// Middleware untuk memeriksa autentikasi JWT
func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 1. Ambil token dari header
		tokenString := c.Get("token")

		if tokenString == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(web.WebResponse{
				Status:  false,
				Message: "Gagal",
				Errors:  "Token tidak ditemukan",
			})
		}

		// 2. Validasi token
		claims, err := helpers.ValidateToken(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(web.WebResponse{
				Status:  false,
				Message: "Gagal",
				Errors:  "Token tidak valid: " + err.Error(),
			})
		}

		// 3. Simpan data user ke 'Context Locals' agar bisa diakses handler
		c.Locals("user_id", claims.UserID)
		c.Locals("role", claims.Role)

		// 4. Lanjutkan ke handler berikutnya
		return c.Next()
	}
}