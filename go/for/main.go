package main

import "fmt"

func main() {
		
	m := make(map[string]int32)

	v, ok := m["yes"]
	if v == 0 {
		fmt.Println("value == 0") // value == 0
		fmt.Println("ok:",ok) // ok: false
	}

	m2 := make(map[string]int32)
	m2["yes"] = 0
	v2, ok2 := m2["yes"]
	if v2 == 0 {
		fmt.Println("value == 0") // value == 0
		fmt.Println("ok:",ok2) // ok: true
	}
}
