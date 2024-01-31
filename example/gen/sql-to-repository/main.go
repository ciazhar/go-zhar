package main

import (
	gen2 "github.com/ciazhar/go-zhar/pkg/gen"
)

func main() {

	err := gen2.SQLToRepository()
	if err != nil {
		panic(err)
	}

}
