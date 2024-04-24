# go:embed

> Go embed adalah fitur yang diperkenalkan pada Go v1.16 yang memungkinkan kita menyertakan file non-Go ke dalam
> executable Go. Saat aplikasi dijalankan, file-file tersebut akan dimuat ke dalam memori, memungkinkan akses cepat dan
> efisien tanpa bergantung pada file eksternal.

Go embed hanya dapat mengakses file yang berada di level folder yang sama atau subfolder-nya. Ini berarti bahwa Go embed
tidak dapat mengambil data dari folder yang berada di level di atasnya. Oleh karena itu, solusi yang dapat digunakan
adalah dengan membuat sebuah FS di dalam folder web untuk mengakses file-file di dalamnya. Setelah itu, FS tersebut
dapat dipanggil di dalam main.go. Berikut adalah cara mendefinisikan FS:

```go
package web

import "embed"

var (
	//go:embed *
	Res embed.FS
)
```

Kemudian FS bisa dipanggil seperti berikut :

```go
package main

import (
	"github.com/ciazhar/go-zhar/examples/go/embed/web"
	"log"
)

func main() {

	file, err := web.Res.ReadFile("static/index.html")
	if err != nil {
		return
	}

	log.Println(string(file))
}

```