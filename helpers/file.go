package helpers

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"time"
)

const (
	TokoImagesPath   = "./public/images/toko"
	ProdukImagesPath = "./public/images/produk"
)

// SaveUploadedFiles menyimpan file ke path yang ditentukan dan mengembalikan nama filenya
func SaveUploadedFiles(files []*multipart.FileHeader, path string) ([]string, error) {
	var filenames []string
	if files == nil {
		return filenames, nil
	}

	for _, file := range files {
		// Buat nama file unik
		filename := fmt.Sprintf("%d-%s", time.Now().UnixNano(), file.Filename)
		filePath := fmt.Sprintf("%s/%s", path, filename)

		// Buka file yang diupload
		src, err := file.Open()
		if err != nil {
			return nil, err
		}
		defer src.Close()

		// Buat file baru di server
		dst, err := os.Create(filePath)
		if err != nil {
			return nil, err
		}
		defer dst.Close()

		// Salin file
		if _, err = io.Copy(dst, src); err != nil {
			return nil, err
		}

		filenames = append(filenames, filename)
	}
	return filenames, nil
}

// DeleteFiles menghapus file dari path yang ditentukan
func DeleteFiles(filenames []string, path string) {
	for _, filename := range filenames {
		if filename == "" {
			continue
		}
		filePath := fmt.Sprintf("%s/%s", path, filename)
		_ = os.Remove(filePath) 
	}
}