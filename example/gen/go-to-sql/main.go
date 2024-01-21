package main

import (
	model2 "github.com/ciazhar/zhar/internal/model"
	gen2 "github.com/ciazhar/zhar/pkg/gen"
)

func main() {

	structList := []gen2.TableDescriber{
		model2.User{},
		model2.Employee{},
	}
	err := gen2.GoToSQL(structList)
	if err != nil {
		panic(err)
	}

}
