# Fiber Clean Architecture

> Part ini akan menjelaskan bagaimana implementasi Clean Architecture menggunakan Fiber Framework di Go.

## Clean Architecture

Clean Architecture adalah sebuah framework pengembangan perangkat lunak yang menekankan pemisahan kode program menjadi
beberapa layer, dimana tiap layer memiliki tanggung jawab masing masing. Layer tersebut yaitu:

- Model, layer untuk mendefinisikan struktur data.
- Repository, layer untuk mendefinisikan komunikasi ke third party lain, seperti database, service dll.
- Service / Use Case, layer untuk mendefinisikan bussiness logic.
- Controller, layer untuk mendefinisikan output aplikasi ke user atau aplikasi lain, seperti REST, gRPC, websocket dll.