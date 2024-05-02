# RabbitMQ: Layanan Antrian Pesan

RabbitMQ adalah perangkat lunak yang memfasilitasi pengiriman pesan antara aplikasi dan service. RabbitMQ berfungsi
sebagai
perantara yang menerima, menyimpan, dan mengirim pesan dari satu aplikasi ke aplikasi lainnya. Penggunaan RabbitMQ
memungkinkan komunikasi asinkron antara aplikasi, meningkatkan scalability dan fault tolerant.

## Pengantar

Contoh kode ini adalah aplikasi sederhana yang menunjukkan penggunaan RabbitMQ untuk mempublikasikan dan mengonsumsi
pesan.

## Penggunaan dan Fitur

### Publish dan Subscribe

RabbitMQ memungkinkan aplikasi untuk mempublikasikan pesan ke antrian dan mengonsumsinya secara asinkron. Fitur ini
berguna untuk:

- Mengirim pesan antara mikro layanan dalam arsitektur berbasis layanan.
- Menerapkan sistem pesan publikasi-subskripsi (pub/sub) untuk pembaruan data real-time.

### TTL (Time-to-Live)

Pada contoh ini, pesan dapat dipublikasikan dengan waktu hidup tertentu menggunakan Time-to-Live (TTL). Fitur ini
bermanfaat untuk:

- Menjadwalkan pemrosesan pesan dalam rentang waktu tertentu.
- Menghapus pesan yang kadaluwarsa dari antrian.

### TestContainers

Pengujian dijalankan dalam lingkungan terisolasi menggunakan TestContainers, yang memungkinkan kami untuk:

- Menjalankan dan mengelola kontainer RabbitMQ selama pengujian.
- Memastikan keandalan aplikasi dalam berinteraksi dengan RabbitMQ dalam berbagai kondisi lingkungan.

## Cara Menjalankan

1. Pastikan RabbitMQ diinstal dan berjalan, atau gunakan kontainer Docker.
2. Jalankan aplikasi dengan menjalankan perintah `go run cmd/main.go`.
3. Akses aplikasi melalui browser atau menggunakan perintah cURL.

Dengan menggunakan RabbitMQ, aplikasi Anda dapat mengirim dan menerima pesan dengan mudah, menjadikannya pilihan yang
kuat untuk membangun sistem yang scalable dan fault tolerant.