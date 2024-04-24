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
