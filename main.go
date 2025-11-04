package main

import (
	"github.com/Debjth19/go-evermos/config"
	"github.com/Debjth19/go-evermos/database"
	"github.com/Debjth19/go-evermos/handler"
	"github.com/Debjth19/go-evermos/repository"
	"github.com/Debjth19/go-evermos/routes"
	"github.com/Debjth19/go-evermos/service"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Inisialisasi Fiber
	app := fiber.New()

	// Koneksi ke Database
	config.ConnectDatabase()

	// Migrasi Database
	database.MigrateDatabase()

	// 1. Repository
	authRepository := repository.NewAuthRepository(config.DB)
	userRepository := repository.NewUserRepository(config.DB)
	alamatRepository := repository.NewAlamatRepository(config.DB)
	tokoRepository := repository.NewTokoRepository(config.DB)
	kategoriRepository := repository.NewKategoriRepository(config.DB)
	produkRepository := repository.NewProdukRepository(config.DB)
	transaksiRepository := repository.NewTransaksiRepository(config.DB)

	// 2. Service
	authService := service.NewAuthService(authRepository)
	userService := service.NewUserService(userRepository, authRepository)
	alamatService := service.NewAlamatService(alamatRepository)
	tokoService := service.NewTokoService(tokoRepository)
	kategoriService := service.NewKategoriService(kategoriRepository)
	produkService := service.NewProdukService(produkRepository, tokoRepository)
	transaksiService := service.NewTransaksiService(config.DB, transaksiRepository, produkRepository, alamatRepository)

	// 3. Handler
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	alamatHandler := handler.NewAlamatHandler(alamatService)
	tokoHandler := handler.NewTokoHandler(tokoService)
	kategoriHandler := handler.NewKategoriHandler(kategoriService)
	produkHandler := handler.NewProdukHandler(produkService)
	transaksiHandler := handler.NewTransaksiHandler(transaksiService)

	// --- Setup Rute ---
	routes.SetupRoutes(app, authHandler, userHandler, alamatHandler, tokoHandler, kategoriHandler, produkHandler, transaksiHandler)
	
	// Rute sederhana untuk tes 
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Selamat Datang di API Evermos",
		})
	})

	// Menjalankan server di port 8000
	err := app.Listen(":8000")
	if err != nil {
		panic("Gagal menjalankan server")
	}
}