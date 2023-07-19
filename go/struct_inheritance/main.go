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

	a3 := A1{B{
		f: 1,
	}}
	a3.SetB(6.0)
	a3.B.SetB(7)                // active
	b3 := a3.B                  // disactive
	b3.SetB(9)                  // disactive
	fmt.Println("a3.f: ", a3.f) // 5.0

	a4 := A2{&B{
		f: 1,
	}}
	a4.SetB(8.0)
	a4.B.SetB(9)                // active
	fmt.Println("a4.f: ", a4.f) // 5.0
}
