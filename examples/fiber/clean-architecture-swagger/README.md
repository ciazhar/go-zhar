# Fiber Clean Architecture & Swagger

> Part ini akan menjelaskan bagaimana implementasi Clean Architecture menggunakan Fiber Framework dan Swagger API doc di
> Go.

## Clean Architecture

Clean Architecture adalah sebuah framework pengembangan perangkat lunak yang menekankan pemisahan kode program menjadi
beberapa layer, dimana tiap layer memiliki tanggung jawab masing masing. Layer tersebut yaitu:

- Model, layer untuk mendefinisikan struktur data.
- Repository, layer untuk mendefinisikan komunikasi ke third party lain, seperti database, service dll.
- Service / Use Case, layer untuk mendefinisikan bussiness logic.
- Controller, layer untuk mendefinisikan output aplikasi ke user atau aplikasi lain, seperti REST, gRPC, websocket dll.

## Swagger

Swagger adalah sebuah framework yang digunakan untuk merancang, mendokumentasikan, dan menguji API. Dengan menggunakan
swagger developer dapat secara visual mendefinisikan spesifikasi API menggunaka format JSON atau YAML. Selain berfungsi
sebagai dokumentasi, swagger juga dapat digunakan untuk testing API secara langsung menggunaakan website.

## Generate Swagger

[Swaggo](https://github.com/swaggo/swag) diperlukan untuk menggenerate Swagger API Doc dari kode program Go. Gunakan
command berikut :

```bash
swag init --parseInternal --dir cmd/,internal/ --output=api/swagger
```

Di Fiber, Swagger API Doc dapat digenerate dengan mendefinisikan komen komen sesuai dengan
format [General API](https://github.com/swaggo/swag?tab=readme-ov-file#general-api-info) pada tiap handler API. Dengan
menggunakan swaggo, swagger dapat di generate menjadi file YAML dan JSON. Gunakan command berikut :

```bash
swag init --parseInternal --dir cmd/,internal/ --output=api/swagger
```

Swaggo akan membaca kode program go yang ada di direktori `cmd` dan `internal` dan kemudian menggenerate Swagger API Doc
di direktori `api/swagger`. Jika berhasil website swagger dapat diakses di http://127.0.0.1:3000/swagger/index.html.