package helpers

import "golang.org/x/crypto/bcrypt"

// Fungsi untuk mengenkripsi password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14) // 14 adalah cost factor
	return string(bytes), err
}

// Fungsi untuk memverifikasi password
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil 
}