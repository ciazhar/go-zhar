package main

import (
	"fmt"
	"unsafe"
)

type BigStruct struct {
	data [1 << 20]int
}

func newBigStruct() *BigStruct {
	var bs BigStruct
	bs.data[10] = 10
	return &bs
}
func main() {
	//bs := newBigStruct()
	//fmt.Println(bs.data[10])
	//// When we're done with bs, it's a good idea to set it to nil to avoid unnecessary memory holding.
	//bs = nil

	//bsa := BigStruct{}
	//
	//// Assign values to the data field
	//for i := 0; i < len(bsa.data); i++ {
	//	bsa.data[i] = i
	//}
	//

	var bsa BigStruct
	bsa2 := BigStruct{}
	bsa3 := new(BigStruct)

	fmt.Printf("Size: %d\n", unsafe.Sizeof(bsa))
	fmt.Printf("Size: %d\n", unsafe.Sizeof(bsa2))
	fmt.Printf("Size: %d\n", unsafe.Sizeof(*bsa3))
	bsa3 = nil
	fmt.Printf("Size: %d\n", unsafe.Sizeof(*bsa3))
}
