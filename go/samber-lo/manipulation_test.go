package lo_test

import (
	"fmt"
	"testing"

	"github.com/samber/lo"
)

func TestToPtr(t *testing.T) {
	str := "hello world"
	ptr := lo.ToPtr(str)
	fmt.Println(ptr, &str)
}

func TestToSlicePtr(t *testing.T) {
	strs := []string{"hello", "world"}
	ptr := lo.ToSlicePtr(strs)
	fmt.Println(ptr)
}

func TestToAnySlices(t *testing.T) {
	data := []int{1, 5, 1}
	elements := lo.ToAnySlice(data)
	fmt.Println(elements)
}

func TestFromAnySlices(t *testing.T) {
	falseData := []any{"foobar", 42}
	elements, ok := lo.FromAnySlice[string](falseData)
	fmt.Println(ok, elements)

	trueData := []any{"foobar", "42"}
	elements, ok = lo.FromAnySlice[string](trueData)
	fmt.Println(ok, elements)
}

func TestEmtpy(t *testing.T) {
	fmt.Println(lo.Empty[int]())
	fmt.Println(lo.Empty[string]())
	fmt.Println(lo.Empty[bool]())
}

func TestCoalesce(t *testing.T) {
	result, ok := lo.Coalesce(0, 1, 2, 3)
	fmt.Println(ok, result)

	res2, ok2 := lo.Coalesce("")
	fmt.Println(ok2, res2)

	var nilString *string
	str := "foobar"
	res3, ok3 := lo.Coalesce[*string](nil, nilString, &str)
	fmt.Println(ok3, res3)
}
