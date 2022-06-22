package lo_test

import (
	"fmt"
	"math"
	"testing"

	"github.com/samber/lo"
)

func TestSubstring(t *testing.T) {
	var sub string
	sub = lo.Substring("hello", 2, 3)
	fmt.Println(sub)

	sub = lo.Substring("hello", -4, 3)
	fmt.Println(sub)

	sub = lo.Substring("hello", -2, math.MaxUint)
	fmt.Println(sub)
}

func TestRuneLength(t *testing.T) {
	var lenght int

	lenght = lo.RuneLength("你好")
	fmt.Println(lenght)

	lenght = len("你好")
	fmt.Println(lenght)
}
