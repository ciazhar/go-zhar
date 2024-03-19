package main

import (
	"github.com/ciazhar/go-zhar/examples/gen/sql-to-repository/model"
	"github.com/ciazhar/go-zhar/pkg/gen"
	"log"
)

func main() {

	structList := []gen.TableDescriber{
		model.User{},
		model.Employee{},
	}
	err := gen.GoToSQL(structList)
	if err != nil {
		log.Fatal(err)
	}

	err = gen.SQLToRepository()
	if err != nil {
		log.Fatal(err)
	}

}
