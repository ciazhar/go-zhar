## bcrypt

bcrypt adalah sebuah library dalam bahasa Go yang digunakan untuk hashing dan memverifikasi password dengan keamanan
yang tinggi. bcrypt didasarkan pada algoritma bcrypt yang didesain khusus untuk menyulitkan proses brute force attack
terhadap password.

### Penggunaan dan Fitur

#### Penggunaan

- **Keamanan Password**: bcrypt digunakan untuk mengamankan password dalam aplikasi dan sistem dengan cara mengubahnya
  menjadi hash yang sulit untuk dipecahkan.
- **Pendaftaran Pengguna**: bcrypt digunakan pada proses pendaftaran pengguna untuk menyimpan password mereka dalam
  bentuk hash, sehingga menghindari penyimpanan password dalam bentuk teks biasa.
- **Otentikasi Pengguna**: bcrypt digunakan untuk memverifikasi password yang dimasukkan oleh pengguna saat proses login
  atau autentikasi.

#### Fitur

- **Kesulitan Brute Force**: bcrypt memungkinkan konfigurasi tingkat kesulitan hashing dengan biaya yang dapat
  disesuaikan (cost), sehingga membuat serangan brute force menjadi lebih sulit dilakukan.
- **Salt Otomatis**: bcrypt secara otomatis menambahkan salt ke dalam hash password, sehingga membuat serangan rainbow
  table menjadi tidak efektif.
- **Keamanan Tinggi**: bcrypt dirancang untuk memberikan keamanan tinggi terhadap serangan kriptografi seperti brute
  force attack, dictionary attack, dan rainbow table attack.

### Contoh Kode (Go)

Berikut adalah contoh penggunaan bcrypt dalam bahasa pemrograman Go:
```go
package main

import (
	"fmt"
	"github.com/ciazhar/go-zhar/pkg/bcrypt"
)

func main() {
	// Password to hash
	password := "mysecretpassword"

	// Hash the password
	hashedPassword, err := bcrypt.HashPassword(password)
	if err != nil {
		fmt.Println("Error hashing password:", err)
		return
	}
	fmt.Println("Hashed password:", hashedPassword)

	// Validate the password
	fmt.Println("Password validation result:", bcrypt.ValidatePassword(password, hashedPassword))
}
```

Dalam contoh ini, kita menggunakan bcrypt untuk menghash dan memverifikasi password. Fungsi `HashPassword` digunakan
untuk menghash password, sedangkan fungsi `ValidatePassword` digunakan untuk memverifikasi apakah password yang
dimasukkan oleh pengguna sesuai dengan hash yang tersimpan.