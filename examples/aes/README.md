## AES (Advanced Encryption Standard)

AES (Advanced Encryption Standard) adalah sebuah algoritma kriptografi simetris yang digunakan untuk mengamankan data
dengan cara mengenkripsi dan mendekripsi pesan. AES adalah standar enkripsi yang paling umum digunakan di dunia saat
ini. Algoritma ini dikenal karena keamanannya yang tinggi dan efisiensinya dalam melakukan proses enkripsi dan dekripsi.

### Penggunaan dan Fitur

#### Penggunaan

- **Pengamanan Data**: AES digunakan untuk mengamankan data dalam berbagai aplikasi dan sistem, termasuk komunikasi
  online, penyimpanan data, dan autentikasi pengguna.
- **Kriptografi File**: AES dapat digunakan untuk mengenkripsi dan mendekripsi file, sehingga data dalam file tersebut
  tidak dapat diakses tanpa kunci yang sesuai.

#### Fitur

- **Keamanan Tinggi**: AES menggunakan kunci simetris yang panjangnya dapat mencapai 128, 192, atau 256 bit, sehingga
  memberikan tingkat keamanan yang tinggi terhadap serangan kriptografi.
- **Efisiensi**: Algoritma AES dirancang untuk mengenkripsi dan mendekripsi data dengan cepat dan efisien, sehingga
  cocok untuk digunakan dalam aplikasi dengan kebutuhan kinerja yang tinggi.
- **Fleksibilitas**: AES dapat diimplementasikan dalam berbagai platform dan bahasa pemrograman, sehingga mudah
  diintegrasikan ke dalam berbagai sistem dan aplikasi.

### Contoh Kode (Go)

Berikut adalah contoh penggunaan AES dalam bahasa pemrograman Go:

```go
package main

import (
	"fmt"
	"github.com/ciazhar/go-zhar/pkg/aes"
)

func main() {
	key := aes.GenerateKey()
	fmt.Println("Key:", key)

	plaintext := "Hello, AES in Go!"

	// Encrypt
	ciphertext := aes.Encrypt(plaintext, key)
	fmt.Println("Encrypted:", ciphertext)

	// Decrypt
	decryptedText := aes.Decrypt(ciphertext, key)
	fmt.Println("Decrypted:", decryptedText)
}

```

Dalam contoh ini, kita menggunakan AES untuk mengenkripsi dan mendekripsi pesan dengan menggunakan kunci yang dihasilkan
secara acak. Pesan awal "Hello, AES in Go!" dienkripsi menjadi ciphertext, lalu didekripsi kembali menjadi plaintext.