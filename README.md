# Mini Proyek Go API Evermos

Ini adalah implementasi API untuk virtual internship Evermos menggunakan Go (Fiber), GORM, dan MySQL, dengan menerapkan Clean Architecture.

## 🚀 Cara Menjalankan

1.  Clone repositori ini.
2.  Pastikan Anda memiliki XAMPP/MySQL yang berjalan.
3.  Buat database di MySQL/MariaDB: `CREATE DATABASE go_evermos;`
4.  Buat file `.env` di root proyek dan isi variabel berikut (sesuaikan dengan setup Anda):
    ```env
    DB_HOST=127.0.0.1
    DB_PORT=3306
    DB_USER=root
    DB_PASS=isi_password_mysql_anda (kosongkan jika pakai XAMPP)
    DB_NAME=go_evermos
    JWT_SECRET=buat_secret_jwt_panjang_dan_acak_anda_sendiri_disini
    ```
5.  Jalankan `go mod tidy` untuk menginstal semua dependensi.
6.  Jalankan server: go run main.go
7.  Server akan berjalan di `http://localhost:8000`.

## 🧪 Pengujian

Semua pengujian endpoint dapat dilakukan menggunakan file `Rakamin Evermos Virtual Internship.postman_collection.json` dari soal.
