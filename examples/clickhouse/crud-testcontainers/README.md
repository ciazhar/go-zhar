# ClickHouse CRUD Testcontainers

## Pengantar

Kode ini adalah contoh implementasi repository untuk menggunakan ClickHouse sebagai database dan Testcontainers untuk
mengelola kontainer ClickHouse dalam unit test. Dalam dokumen ini, kita akan membahas apa itu ClickHouse, apa itu
Testcontainers, dan beberapa contoh penggunaannya.

## Apa itu ClickHouse?

ClickHouse adalah database analitik berkinerja tinggi yang dikembangkan oleh Yandex. Ini dirancang khusus untuk analisis
data besar dengan kecepatan tinggi. ClickHouse mendukung pengolahan SQL paralel, kompresi data, pengindeksan kolom, dan
banyak lagi fitur untuk mempercepat kueri analitik pada data besar.

## Apa itu Testcontainers?

Testcontainers adalah library Go yang memungkinkan Anda untuk dengan mudah membuat dan mengelola kontainer Docker dalam
unit test. Ini memungkinkan testing yang konsisten dan dapat diulang, serta isolasi environment testing.

## Penggunaan dan Fitur

3. **Membuat Event**: Fungsi `CreateEvent` digunakan untuk membuat sebuah event baru dalam tabel `events`.
4. **Mendapatkan Event**: Fungsi `GetEvent` digunakan untuk mendapatkan event berdasarkan `event_id`
   dan `injection_time`.
5. **Mendapatkan Daftar Event**: Fungsi `GetEvents` digunakan untuk mendapatkan daftar event berdasarkan jenis (`types`)
   dan/atau penerima (`rcpTo`).
6. **Mendapatkan Daftar Event dengan Cursor**: Fungsi `GetEventsCursor` digunakan untuk mendapatkan daftar event dengan
   menggunakan cursor pagination.
7. **Mendapatkan Agregasi Harian**: Fungsi `GetAggregateDaily` digunakan untuk mendapatkan agregasi data harian dari
   tabel `events`.
8. **Mendapatkan Agregasi Per Jam**: Fungsi `GetAggregateHourly` digunakan untuk mendapatkan agregasi data per jam dari
   tabel `events`.

Dengan menggunakan ClickHouse dan Testcontainers, Anda dapat dengan mudah mengembangkan, menguji, dan menerapkan
aplikasi yang mampu menangani analisis data besar dengan cepat dan efisien.