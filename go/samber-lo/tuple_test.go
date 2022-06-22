package lo_test

import (
	"fmt"
	"testing"

	"github.com/samber/lo"
)

func example() (string, int) { return "y", 2 }
func TestT2(t *testing.T) {
	tuple1 := lo.T2("x", 1)
	fmt.Println(tuple1)

	tuple2 := lo.T2(example())
	fmt.Println(tuple2)

}

func TestUnpack2(t *testing.T) {
	tuple1 := lo.T2("x", 1)
	r1, r2 := lo.Unpack2(tuple1)
	fmt.Println(r1, r2)

}

func TestZip2(t *testing.T) {
	tp2 := lo.Zip2([]string{"a", "b"}, []int{1, 2})
	fmt.Println(tp2)

	tp4 := lo.Zip2([]string{"a", "b", "c"}, []int{1, 2, 3, 4})
	fmt.Println(tp4)
}

func TestUnzip2(t *testing.T) {
	tp2 := lo.Zip2([]string{"a", "b"}, []int{1, 2})
	fmt.Println(tp2)

	a, b := lo.Unzip2(tp2)
	fmt.Println(a, b)
}
