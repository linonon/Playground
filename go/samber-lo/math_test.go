package lo_test

import (
	"fmt"
	"testing"

	"github.com/samber/lo"
)

func TestRangeAndRangeFromAndRangeWithSteps(t *testing.T) {
	var r []int
	var rf []float64

	r = lo.Range(4)
	fmt.Println(r)

	r = lo.Range(-4)
	fmt.Println(r)

	r = lo.RangeFrom(1, 5)
	fmt.Println(r)

	rf = lo.RangeFrom(1.0, 5)
	fmt.Println(rf)

	r = lo.RangeWithSteps(0, 20, 5)
	fmt.Println(r)

	rf = lo.RangeWithSteps(-1.0, -4.0, -0.5)
	fmt.Println(rf)

	rf = lo.RangeWithSteps(-1.0, -4.0, 0.5)
	fmt.Println(rf)

	r = lo.RangeWithSteps(1, 4, -1)
	fmt.Println(r)

	r = lo.Range(0)
	fmt.Println(r)
}

func TestClamp(t *testing.T) {
	// 如果 value 大於邊界, 則取邊界的值
	v1 := 0
	r1 := lo.Clamp(v1, -10, 10)
	fmt.Println("v1:", v1, "r1:", r1)

	v2 := -42
	r2 := lo.Clamp(v2, -10, 10)
	fmt.Println("v2:", v2, "r2:", r2)

	v3 := 42
	r3 := lo.Clamp(v3, -10, 10)
	fmt.Println("v3:", v3, "r3:", r3)
}

func TestSumBy(t *testing.T) {
	strs := []string{"foo", "bar"}
	callback := func(item string) int {
		return len(item)
	}

	sum := lo.SumBy(strs, callback)
	fmt.Println(sum)
}
