# Clean Architecture
> Clean Architecture push us to separate stable business rules (higher-level abstractions) from volatile technical details (lower-level details), defining clear boundaries. The main building block is the Dependency Rule : source code dependencies must point only inward, toward higher-level policies. (Robert C. Martin 2017)


Disini saya cuman akan menjelaskan secara singkat. Lebih lenkapnya bisa baca [artike nya mas iman](https://hackernoon.com/golang-clean-archithecture-efd6d7c43047). Ambil contoh ada API untuk memasukkan data artikel ke database sebagai berikut :
```go
package main

import (
    "github.com/astaxie/beego/orm"
    "github.com/gin-gonic/gin"
)

type Article struct {
    ID string   `json:"id"`
    Name string `json:"name"`
    Slug string `json:"slug"`
}

func Store(c *gin.Context)  {
    
    //request body
    var payload Article
    if err := c.Bind(&payload); err != nil {
        c.JSON(http.StatusBadRequest, log.Warn(err))
        return
    }
    
    //slugify article name
    payload.Slug=slugify(payload.name) 

    //insert db
    if _, err := orm.NewOrm().Insert(payload); err != nil {
        c.JSON(http.StatusInternalServerError, log.Error(err))
        return
    }

    c.JSON(http.StatusOK, payload)
}
```
Untuk versi full nya bisa dilihat di [sini](plain) 

### Pembagian Layer
Nanti api tersebut akan dibagi ke beberapa layer yaitu :
#### Model
Layer ini merupakan layer yang menyimpan model yang dipakai pada domain lainnya. Layer ini dapat diakses oleh semua layer dan oleh semua domain.
#### Repository
Layer ini merupakan layer yang menyimpan database handler. Querying, Inserting, Deleting akan dilakukan pada layer ini. Tidak ada business logic disini. Yang ada hanya fungsi standard untuk input-output dari datastore.
Layer ini memiliki tugas utama yakni menentukan datastore apa yang di gunakan. Teman-teman boleh memilih sesuai kepada kebutuhan, mungkin RDBMS (Mysql,PostgreSql, dsb) atau NoSql (Mongodb,CouchDB dsb).
Jika menggunakan arsitektur microservice, maka layer ini akan bertugas sebagai penghubung kepada service lain. Layer ini akan terikat dan bergantung pada datastore yang digunakan.
#### Use Case
Layer ini merupakan layer yang akan bertugas sebagai pengontrol, yakni menangangi business logic pada setiap domain. Layer ini juga bertugas memilih repository apa yang akan digunakan, dan domain ini bisa memiliki lebih dari satu repository layer.
Tugas utama terbesar dari layer ini, yaitu menjadi penghubung antara datastore (repository layer) dengan delivery layer. Sehingga, layer ini juga bertanggung jawab atas kevalidan data, jika sesuatu terjadi data yang tidak valid pada repository atau delivery, maka layer ini yang pertama kali disalahkan.
Layer ini benar benar harus berisi business logic, contohnya: penjumlahan, total masukan, atau membentuk response yang merupakan gabungan dari beberapa repository/model. Layer ini bergantung pada repository layer. Jadi jika terjadi perubahan di repository secara besar-besaran tentu saja mempengaruhi layer ini.
### Controller
Layer ini merupakan layer yang akan bertugas sebagai presenter atau menjadi output dari aplikasi. Layer ini bertugas menentukan metode penyampaian yang dipakai, bisa dengan Rest API, HTML, gRPC, File dsb.
Tugas lain dari layer ini, menjadi dinding penghubung antara user dan sistem. Menerima segala input dan validasi input sesuai standar yang digunakan.
Untuk contoh project yang saya gunakan, saya memilih Rest API sebagai delivery layernya. Sehingga, komunikasi antara client/user terhadap sistem dilakukan melalui REST API
Dia berisi protocol delivery ke user

### Kenapa Dibagi ?
- Mempermudah testing
- Koding lebih terstruktur
- Low Bug 

Sehingga nanti hasilnya seperti [disini](clean-arch-example)
