package main

import (
	"fmt"
	"unsafe"
)

func main() {
	p := new(person)

	pName := (*string)(unsafe.Pointer(p))
	*pName = "linonon"

	pAge := (*int)(unsafe.Pointer(uintptr(unsafe.Pointer(p)) + unsafe.Offsetof(p.Age)))
	*pAge = 20
	fmt.Println(*p)
}

type person struct {
	Name string
	Age  int
}
