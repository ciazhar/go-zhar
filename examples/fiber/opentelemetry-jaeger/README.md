# Distributed Tracing With Open Telemetry and Jaeger in Fiber

## Pengenalan
Part ini adalah implementasi microservice menggunakan framework Fiber di Go, dengan integrasi OpenTelemetry dan Jaeger untuk observabilitas. Proyek ini terdiri dari dua service utama: `order` dan `user`.

### Apa itu OpenTelemetry?
OpenTelemetry adalah tools, API, dan SDK yang digunakan untuk mengumpulkan, memproses, dan mengekspor data telemetri (tracing, metrics, dan log) dari software untuk membantu dalam observabilitas sistem.

### Apa itu Jaeger?
Jaeger adalah sistem open-source untuk pelacakan transaksi terdistribusi. Ini digunakan untuk memantau dan memecahkan masalah dalam arsitektur microservices, mengukur kinerja, dan menganalisis ketergantungan service.

## Struktur Proyek
Proyek ini dibagi menjadi beberapa direktori utama:
- `cmd/order`: Berisi kode aplikasi untuk layanan `order`.
- `cmd/user`: Berisi kode aplikasi untuk layanan `user`.
- `internal/order`: Berisi logika bisnis dan implementasi dari layanan `order`.
- `internal/user`: Berisi logika bisnis dan implementasi dari layanan `user`.
- `pkg`: Berisi paket umum yang digunakan oleh berbagai layanan, seperti logger, konfigurasi environment, dan integrasi OpenTelemetry-Jaeger.

## Fitur
- **Order Service**:
    - Menambahkan order baru.
    - Mengambil order berdasarkan ID.
    - Mengambil semua order.
    - Menghapus order.
    - Memperbarui order.

- **User Service**:
    - Menambahkan user baru.
    - Mengambil user berdasarkan username.
    - Mengambil semua user.
    - Menghapus user.
    - Memperbarui user.

- **Observabilitas**:
    - Integrasi dengan OpenTelemetry untuk pelacakan distribusi.
    - Penggunaan Jaeger untuk mengumpulkan dan menganalisis tracing.

## Use Case
1. **Manajemen Order**:
    - Service `order` memungkinkan untuk menambahkan, mengambil, memperbarui, dan menghapus order. Order dapat dikaitkan dengan user yang terdaftar di service `user`.

2. **Manajemen User**:
    - Service `user` menyediakan endpoint untuk menambah, mengambil, memperbarui, dan menghapus user. Data user digunakan oleh service `order` untuk mengaitkan order dengan user tertentu.

3. **Monitoring dan Observabilitas**:
    - Dengan mengintegrasikan OpenTelemetry dan Jaeger, setiap request yang masuk ke service `order` dan `user` akan dilacak. Ini membantu dalam mengidentifikasi masalah performa dan ketergantungan antar service.

## Instalasi dan Penggunaan
- Clone repository ini.
- Pastikan Anda memiliki Go, Docker, dan Docker Compose terinstal.
- Jalankan docker-compose up -d untuk menjalankan Jaeger. Jaeger UI dapat diakses di http://localhost:16686/.
- Jalankan layanan order dan user menggunakan perintah go run cmd/order/main.go dan go run cmd/user/main.go.
= Layanan sekarang dapat diakses melalui endpoint yang telah dikonfigurasi, misalnya http://localhost:3000/orders.