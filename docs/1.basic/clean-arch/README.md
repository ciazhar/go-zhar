# Clean Architecture
> Clean Architecture push us to separate stable business rules (higher-level abstractions) from volatile technical details (lower-level details), defining clear boundaries. The main building block is the Dependency Rule : source code dependencies must point only inward, toward higher-level policies. (Robert C. Martin 2017)

Ambil contoh ada API untuk memasukkan data artikel ke database sebagai berikut :
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

Nanti api tersebut akan dibagi ke beberapa layer yaitu :
- Model

Dia define struktur data

- Repository

Dia berisi query db, http client

- Use Case

Dia berisi logic

- Controller

Dia berisi protocol delivery ke user



Kenapa Dibagi ?
- Mempermudah testing
- Koding lebih terstruktur
- Low Bug 

Sehingga nanti hasilnya seperti [disini](clean-arch-example)
