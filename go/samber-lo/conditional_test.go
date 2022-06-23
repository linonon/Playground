package lo_test

import (
	"fmt"
	"math"
	"testing"

	"github.com/samber/lo"
)

func TestTernary(t *testing.T) {
	var res string

	res = lo.Ternary(true, "a", "b")
	fmt.Println(res)

	res = lo.Ternary(false, "a", "b")
	fmt.Println(res)
}

func TestIfElseIfElse(t *testing.T) {
	var x, y int

	x, y = 1, 3
	res := lo.If(x < y, 1).
		ElseIf(x < y-1, 2).
		Else(3)
	fmt.Println(res)

	x, y = 2, 16
	res2 := lo.If(int(math.Pow(float64(x), 5)) == y, "x^5 == y").
		ElseIf(int(math.Pow(float64(x), 4)) == y, "x^4 == y").
		Else("x^4 != y && x^5 !=y")
	fmt.Println(res2)

	res3 := lo.IfF(true, func() int { return 1 }).
		ElseIfF(false, func() int { return 2 }).
		ElseF(func() int { return 3 })
	fmt.Println(res3)
}

func TestSwitchCaseDefalut(t *testing.T) {
	res := lo.Switch[int, string](1).
		Case(1, "1").
		Case(2, "2").
		Default("3")
	fmt.Println(res)

	res2 := lo.Switch[int, string](1).
		CaseF(1, func() string { return "1" }).
		CaseF(2, func() string { return "2" }).
		DefaultF(func() string { return "nil" })
	fmt.Println(res2)

}
