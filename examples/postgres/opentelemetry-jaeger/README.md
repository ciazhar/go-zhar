# CRUD dengan PostgreSQL menggunakan sqlc

## Deskripsi

Proyek ini merupakan contoh implementasi operasi CRUD (Create, Read, Update, Delete) menggunakan PostgreSQL dan sqlc di
Go. CRUD adalah konsep dasar dalam pengembangan perangkat lunak yang memungkinkan pengguna untuk membuat, membaca,
memperbarui, dan menghapus data dalam basis data.

## PostgreSQL

PostgreSQL adalah sistem manajemen basis data (DBMS) yang open-source dan kuat. Ini menyediakan banyak fitur canggih
termasuk dukungan untuk jenis data yang beragam, fungsi tingkat tinggi, dan fitur skalabilitas yang kuat.

## CRUD

CRUD adalah singkatan dari Create, Read, Update, dan Delete. Ini adalah operasi dasar yang digunakan dalam pengembangan
perangkat lunak untuk memanipulasi data dalam sebuah sistem. Operasi-operasi ini meliputi:

- **Create**: Membuat entitas baru dalam basis data.
- **Read**: Membaca entitas dari basis data.
- **Update**: Memperbarui entitas yang ada dalam basis data.
- **Delete**: Menghapus entitas dari basis data.

## sqlc

sqlc adalah alat yang memungkinkan pengguna untuk menulis query SQL dalam file .sql dan kemudian menghasilkan kode Go
yang menggunakan query-query ini sebagai metode dalam pemrograman Go. Hal ini memungkinkan pemisahan yang jelas antara
logika aplikasi dan query SQL-nya, serta menyediakan fitur autentikasi terhadap kesalahan saat waktu kompilasi.

## Fitur dan Usecase

### Fitur

1. **CreateProduct**: Membuat produk baru dengan nama dan harga yang ditentukan.
2. **GetProducts**: Mendapatkan daftar produk berdasarkan kriteria pencarian tertentu dengan dukungan untuk pengurutan
   dan pembagian halaman.
3. **GetProductsCursor**: Mendapatkan daftar produk dengan menggunakan cursor untuk pembagian halaman yang lebih
   efisien.
4. **UpdateProductPrice**: Memperbarui harga produk yang ada.
5. **DeleteProduct**: Menghapus produk dari basis data.

### Usecase

- Sebagai admin toko online, saya ingin menambahkan produk baru ke dalam sistem.
- Sebagai pengguna, saya ingin melihat daftar produk yang tersedia dalam toko online dengan kemampuan untuk mencari dan
  mengurutkannya.
- Sebagai pengguna, saya ingin melihat daftar produk dalam bentuk halaman-halaman yang bisa saya navigasikan menggunakan
  tombol next dan prev.

## Cara Menjalankan Proyek

1. Pastikan Anda telah mengonfigurasi koneksi basis data PostgreSQL pada file `config.json`.
2. Jalankan aplikasi dengan menjalankan perintah `go run main.go`.

Dengan menggunakan sqlc dan Go, proyek ini memberikan cara yang efisien dan andal untuk melakukan operasi CRUD pada
basis data PostgreSQL.