## README

### Overview

Repository ini akan membahas implementasi Unit Test menggunakan mekanisme Mock di Clean Architecture layer Service dan
Controller menggunakan library `testify` dan `mockery`. Layer ini membantu memisahkan logika bisnis dari logika akses
data, membuat kode lebih terstruktur, mudah diuji, dan dikelola.

### Struktur Repository

- `repository`: Mengelola akses data.
- `service`: Mengelola logika bisnis.
- `controller`: Mengelola HTTP handler.

### Dependensi

- **Testify**: Library untuk assertion dan mock dalam pengujian unit.
- **Mockery**: Generator kode untuk membuat mock object dari interface yang ada. Dengan Mockery, kita bisa
  membuat mock object yang dapat digunakan untuk testing, memungkinkan kita untuk mengisolasi dan menguji bagian
  tertentu dari kode tanpa tergantung pada implementasi konkret.

### Installasi Mockery
```sh
  go install github.com/vektra/mockery/v2/.../
```

#### Cara Menggunakan Mockery

1. Buat interface yang ingin di-mock.
2. Generate mock menggunakan mockery:
   ```sh
   mockery  --output ./internal/mocks --dir ./internal --all
   ```
3. Gunakan mock dalam pengujian unit, seperti contoh pada `service_test.go` dan `controller_test.go`.