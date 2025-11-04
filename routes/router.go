package routes

import (
	"github.com/Debjth19/go-evermos/handler"
	"github.com/Debjth19/go-evermos/middleware"

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes mengkonfigurasi semua rute API
func SetupRoutes(
	app *fiber.App,
	authHandler handler.AuthHandler,
	userHandler handler.UserHandler,
	alamatHandler handler.AlamatHandler,
	tokoHandler handler.TokoHandler, 
	kategoriHandler handler.KategoriHandler,
	produkHandler handler.ProdukHandler,
	transaksiHandler handler.TransaksiHandler,
) {
	api := app.Group("/api/v1")

	// Rute untuk Autentikasi 
	auth := api.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)

	// Rute untuk User
	user := api.Group("/user", middleware.AuthMiddleware()) // Middleware di sini
	user.Get("/", userHandler.GetProfile)
	user.Put("/", userHandler.UpdateProfile)

	// Rute untuk Alamat
	alamat := user.Group("/alamat") 
	alamat.Post("/", alamatHandler.CreateAlamat)
	alamat.Get("/", alamatHandler.GetAllAlamat)
	alamat.Get("/:id", alamatHandler.GetAlamatByID)
	alamat.Put("/:id", alamatHandler.UpdateAlamat)
	alamat.Delete("/:id", alamatHandler.DeleteAlamat)

	// Rute untuk Toko
	toko := api.Group("/toko")
	
	// Rute yang perlu autentikasi
	toko.Get("/my", middleware.AuthMiddleware(), tokoHandler.GetMyToko)
	toko.Put("/:id_toko", middleware.AuthMiddleware(), tokoHandler.UpdateToko)
	
	// Rute publik 
	toko.Get("/", tokoHandler.GetAllToko) // -> /api/v1/toko
	toko.Get("/:id_toko", tokoHandler.GetTokoByID) // -> /api/v1/toko/:id_toko

	// Rute untuk Kategori (Perlu Token & Role Admin)
	category := api.Group("/category", middleware.AuthMiddleware(), middleware.AdminMiddleware())
	
	// Terapkan middleware ke semua rute di grup ini
	category.Post("/", kategoriHandler.CreateKategori)
	category.Get("/", kategoriHandler.GetAllKategori)
	category.Get("/:id", kategoriHandler.GetKategoriByID)
	category.Put("/:id", kategoriHandler.UpdateKategori)
	category.Delete("/:id", kategoriHandler.DeleteKategori)

	// Rute untuk Produk
	product := api.Group("/product")

	// Rute yang perlu autentikasi
	product.Post("/", middleware.AuthMiddleware(), produkHandler.CreateProduk)
	product.Put("/:id", middleware.AuthMiddleware(), produkHandler.UpdateProduk)
	product.Delete("/:id", middleware.AuthMiddleware(), produkHandler.DeleteProduk)

	// Rute publik 
	product.Get("/", produkHandler.GetAllProduk)
	product.Get("/:id", produkHandler.GetProdukByID)

	// Rute untuk Transaksi (Perlu Autentikasi)
	trx := api.Group("/trx", middleware.AuthMiddleware())
	trx.Post("/", transaksiHandler.CreateTransaksi)
	trx.Get("/", transaksiHandler.GetMyTransactions)
	trx.Get("/:id", transaksiHandler.GetMyTransactionByID)


}