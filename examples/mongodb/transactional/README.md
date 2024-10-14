# MongoDB Transactional

> Part ini akan menjelaskan bagaimana implementasi transaction pada MongoDB di Go. Untuk detail mengenai transaction di
> MongoDB bisa di cek di [laman berikut](https://www.mongodb.com/docs/manual/core/transactions-in-applications/)

## Transaction

## MongoDB Transaction

Terdapat 2 cara untuk melakukan transaction yaitu dengan menggunakan `WithTransaction` dan `StartTransaction`. Perbedaan
yang mendasar adalah `WithTransaction` merupakan implementasi transaction yang simple, semua error handling sudah di
tangani oleh `WithTransaction` dan kita hanya perlu mendeklarasikan database interaction di dalam transaction function.
Sedangkan `StartTransaction` lebih customable, kita bisa mengcustom bagaimana cara menghandle error, retry mechanism,
step mana yang akan trigger abort transaction dll. Untuk implementasi `WithTransaction` bisa dilihat
di `PurchaseWithAutomaticTransaction` sedangkan implementasi `StartTransaction` berada
di `PurchaseWithManualTransaction`.
Transaction di MongoDB sendiri membutuhkan MongoDB dijalankan sebagai cluster.

## User Story

Aplikasi ini merupakan aplikasi yang mendata buku beserta dengan stok nya. Data buku akan disimpan di collection books.
Kita bisa membeli buku tersebut dan kemudian stok nya akan berkurang. Selain itu juga akan terdapat log transaksi
terkait berapa jumlah buku yang dibeli. Transaksi tersebut disimpan di collection transaction.

## Run MongoDB Cluster

```bash
cd ../../../deployments/mongodb && docker-compose up
```

## Run Server Service

Untuk menjalankan server service, gunakan perintah berikut.

```bash
go run cmd/server.go
````

