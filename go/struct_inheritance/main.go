package main

import "fmt"

type A1 struct {
	B
}

type A2 struct {
	*B
}

type B struct {
	f float64
}

func (b *B) SetB(f float64) {
	b.f = f
}

func main() {
	b1 := B{1.0}
	a1 := A1{B: b1}
	a1.SetB(2.0)
	b1.SetB(3)                  // disactive
	fmt.Println("a1.f: ", a1.f) // 2.0

	b2 := B{1.0}
	a2 := A2{B: &b2}
	a2.SetB(4.0)
	b2.SetB(5)                  // active
	fmt.Println("a2.f: ", a2.f) // 5.0
}
