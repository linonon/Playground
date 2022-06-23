package lo_test

import (
	"fmt"
	"testing"

	"github.com/samber/lo"
)

func TestEvery(t *testing.T) {
	data := []int{0, 1, 2, 3, 4, 5, 6, 7}
	subset := []int{1, 2, 3, 4, 5, 6}
	ok := lo.Every(data, subset)
	fmt.Println(ok)

	subset2 := []int{4, 5, 6, 7, 8}
	ok2 := lo.Every(data, subset2)
	fmt.Println(ok2)
}

func TestEveryBy(t *testing.T) {
	data := []int{1, 2, 3, 4}
	callback := func(x int) bool { return x < 5 }

	ok := lo.EveryBy(data, callback)
	fmt.Println(ok)
}

func TestSome(t *testing.T) {
	data := []int{1, 2, 3, 4, 5}

	subset := []int{0, 2}
	ok := lo.Some(data, subset)
	fmt.Println(ok)

	subset2 := []int{-1, 6}
	ok2 := lo.Some(data, subset2)
	fmt.Println(ok2)
}

func TestSomeBy(t *testing.T) {
	data := []int{1, 2, 3, 4, 5}
	callback := func(x int) bool { return x < 3 }

	ok := lo.SomeBy(data, callback)
	fmt.Println(ok)
}

func TestNone(t *testing.T) {
	data := []int{1, 2, 3, 4, 5}

	subset := []int{0, 2}
	ok := lo.None(data, subset)
	fmt.Println(ok)

	subset2 := []int{-1, 6}
	ok2 := lo.None(data, subset2)
	fmt.Println(ok2)
}

func TestNoneBy(t *testing.T) {
	data := []int{1, 2, 3, 4, 5}
	lessThan3 := func(x int) bool { return x < 3 }

	ok := lo.NoneBy(data, lessThan3)
	fmt.Println(ok)

	lessThan0 := func(x int) bool { return x < 0 }
	ok2 := lo.NoneBy(data, lessThan0)
	fmt.Println(ok2)
}

func TestIntersect(t *testing.T) {
	data := []int{1, 2, 3, 4, 5}

	sub := []int{1, 4}
	r1 := lo.Intersect(data, sub)
	fmt.Println(r1)

	sub2 := []int{1, 6}
	r2 := lo.Intersect(data, sub2)
	fmt.Println(r2)

	sub3 := []int{0, 6}
	r3 := lo.Intersect(data, sub3)
	fmt.Println(r3)
}

func TestDifference(t *testing.T) {
	l1 := []int{0, 1, 2, 3, 4, 5}
	l2 := []int{0, 2, 6, 4}

	l, r := lo.Difference(l1, l2)
	fmt.Println(l, r)
}

func TestUnion(t *testing.T) {
	l1 := []int{0, 1, 2, 3, 4, 5}
	l2 := []int{0, 2, 6, 4}

	union := lo.Union(l1, l2)
	fmt.Println(union)
}
