# Golden File Testing

Golden File Testing adalah sebuah teknik testing di Go, dimana untuk memvalidasi apakah response hasil suatu function
sama dengan data yang sudah di ekspektasikan yang disimpan dalam file golden.

Kita disini memiliki function `GenerateOutput()` yang mengoutputkan string sesuai dengan input pada parameter yang
diberikan. Kita akan memastikan output tetap konsisten. Kita akan membuat testing untuk fungsi ini menggunakan golden
file testing.

```go
// main.go
package main

import (
	"fmt"
)

func GenerateOutput(input string) string {
	// Some logic to generate output based on input
	return "Output: " + input
}

func main() {
	output := GenerateOutput("example")
	fmt.Println(output)
}

```

```go
// main_test.go
package main

import (
	"os"
	"testing"
)

func TestGenerateOutput(t *testing.T) {
	// Input
	input := "example"

	// Call the function
	output := GenerateOutput(input)

	// Read the golden file
	golden, err := os.ReadFile("testdata/data.golden")
	if err != nil {
		t.Fatalf("unable to read golden file: %v", err)
	}

	// Compare output with golden file content
	if string(golden) != output {
		t.Errorf("output does not match golden file:\nExpected: %s\nGot: %s", golden, output)
	}
}

```

Pada kode di atas:

- `main.go` berisi function `GenerateOutput()` yang ingin kita uji.
- `main_test.go` berisi test function `TestGenerateOutput()` yang membandingkan output dari fungsi `GenerateOutput()`
  dengan konten dari file golden yang berlokasi di `testdata/output.golden`.

Selanjutnya, Kita perlu membuat file golden `output.golden` di dalam direktori bernama `testdata`. Isi dari file ini
harus berupa output yang diharapkan dari `GenerateOutput()` untuk input yang diberikan.

Berikut struktur direktorinya:

```go
project/
│
├── main.go
├── main_test.go
└── testdata/
└── golden_output.txt
```

Pastikan untuk mengisi `output.golden` dengan output yang diharapkan. Ketika pengujian dijalankan, hasil output aktual
dari `GenerateOutput()` akan dibandingkan dengan konten dari output.golde. Jika sesuai, pengujian akan lulus. Jika
tidak, pengujian akan gagal, menandakan bahwa output telah berubah secara tidak terduga.