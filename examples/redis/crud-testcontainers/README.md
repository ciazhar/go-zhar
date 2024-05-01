# CRUD Testcontainers with Redis

## Pendahuluan

Repository ini merupakan contoh penggunaan Redis dalam aplikasi CRUD menggunakan testcontainers sebagai pengujian.

## Apa itu Redis?

Redis adalah sistem penyimpanan data open-source yang memungkinkan data disimpan pada memori utama. Ini dapat digunakan
sebagai database, cache, dan broker pesan. Redis menyediakan berbagai struktur data seperti string, hash, lists, sets,
sorted sets dengan query yang sangat kuat.

## Apa itu Testcontainers?

Testcontainers adalah suatu framework yang memungkinkan untuk membuat kontainer Docker pada pengujian, yang berguna
untuk mengisolasi pengujian Anda dari lingkungan lokal Anda, sehingga pengujian dapat diulang dengan konsisten pada
berbagai lingkungan.

## Use Case

- `Get() (string, error)`: Mengambil nilai dari Redis berdasarkan kunci yang diberikan.
- `Set(value string, expiration time.Duration) error`: Menyimpan nilai dalam Redis dengan kunci dan masa berlaku yang
  ditentukan.
- `GetHash(field string) (string, error)`: Mengambil nilai dari hash dalam Redis berdasarkan field yang diberikan.
- `SetHash(field string, value string) error`: Menyimpan nilai dalam hash Redis dengan field dan nilai yang diberikan.
- `SetHashTTL(field string, value string, ttl time.Duration) error`: Menyimpan nilai dalam hash Redis dengan field,
  nilai, dan TTL (Time To Live) yang ditentukan.
- `DeleteHash(field string) error`: Menghapus nilai dari hash Redis berdasarkan field yang diberikan.

## Cara Menjalankan

Pastikan Docker sudah terpasang di sistem Anda.
Jalankan perintah `go test ./internal/repository` untuk menjalankan pengujian menggunakan Redis container.