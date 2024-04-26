# Go Validator

Package GoValidator ini menyediakan serangkaian tools untuk memvalidasi struct dalam Go.

# Penggunaan

Berikut adalah contoh dasar penggunaan GoValidator:

```go
package main

import (
	"github.com/ciazhar/go-zhar/pkg/logger"
	"github.com/ciazhar/go-zhar/pkg/validator"
)

type User struct {
	Username string `validate:"required"`
	Tagline  string `validate:"required,lt=10"`
}

func main() {
	// Inisialisasi logger
	log := logger.Init()

	// Inisialisasi validator untuk Bahasa Inggris
	validate := validator.New("en", log)

	// Buat instance user
	user := User{
		Username: "Joeybloggs",
		Tagline:  "This tagline is way too long.",
	}

	// Validasi struktur user
	err := validate.ValidateStruct(user)
	if err != nil {
		log.Infof("validateStruct : %v", err)
	}
}
```

# Fitur

- Validasi dalam Beberapa Bahasa: GoValidator mendukung error message dalam beberapa bahasa. Saat ini support Inggris (
  en) dan
  Indonesia (id).
- Validasi Kustom: Anda dapat mendefinisikan rules validasi custom untuk struct dengan mendaftarkan function validasi
  custom.
- Penggantian Terjemahan: Anda dapat mengganti error message validasi default dengan pesan custom untuk tag validasi
  tertentu.

Full example ada di [main.go](https://github.com/ciazhar/go-zhar/blob/master/examples/govalidator/main.go)