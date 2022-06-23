package lo_test

import (
	"fmt"
	"testing"

	"github.com/samber/lo"
)

func TestIndexOf(t *testing.T) {
	data := []int{0, 1, 2, 3, 4, 5}
	v := 2
	found := lo.IndexOf(data, v)
	fmt.Println(found)

	v2 := 6
	notFound := lo.IndexOf(data, v2)
	fmt.Println(notFound)
}

func TestLastIndexOf(t *testing.T) {
	data := []int{0, 1, 2, 2, 3, 3}
	v := 2

	found := lo.LastIndexOf(data, v)
	fmt.Println(found)

	v2 := 6
	notFound := lo.LastIndexOf(data, v2)
	fmt.Println(notFound)
}

func TestFind(t *testing.T) {
	data := []string{"a", "b", "c", "d"}
	callback := func(str string) bool {
		return str == "b"
	}

	str, ok := lo.Find(data, callback)
	fmt.Println(str, ok)
}

func TestFindIndexOf(t *testing.T) {
	data := []string{"a", "b", "a", "b"}
	callback := func(str string) bool {
		return str == "b"
	}

	str, index, ok := lo.FindIndexOf(data, callback)
	fmt.Println(str, index, ok)

	data2 := []string{"foobar"}
	callback2 := func(str string) bool {
		return str == "b"
	}

	str2, index2, ok2 := lo.FindIndexOf(data2, callback2)
	fmt.Println(str2, index2, ok2)
}

func TestFindLastIndexOf(t *testing.T) {
	data := []string{"a", "b", "a", "b"}
	callback := func(str string) bool {
		return str == "b"
	}

	str, index, ok := lo.FindLastIndexOf(data, callback)
	fmt.Println(str, index, ok)

	data2 := []string{"foobar"}
	callback2 := func(str string) bool {
		return str == "b"
	}

	str2, index2, ok2 := lo.FindLastIndexOf(data2, callback2)
	fmt.Println(str2, index2, ok2)
}

func TestMin(t *testing.T) {
	min := lo.Min([]int{1, 2, 3})
	fmt.Println(min)

	min2 := lo.Min([]int{})
	fmt.Println(min2)
}

func TestMinBy(t *testing.T) {
	data := []string{"s1", "string2", "s3"}
	callback := func(item string, min string) bool {
		return len(item) < len(min)
	}
	min := lo.MinBy(data, callback)
	fmt.Println(min)
}

func TestMax(t *testing.T) {
	data := []int{1, 2, 3}
	max1 := lo.Max(data)
	fmt.Println(max1)

	data2 := []int{}
	max2 := lo.Max(data2)
	fmt.Println(max2)
}

func TestMaxBy(t *testing.T) {
	data := []string{"string1", "s2", "string3"}
	callback := func(item string, max string) bool {
		return len(item) > len(max)
	}

	res := lo.MaxBy(data, callback)
	fmt.Println(res)
}

func TestLast(t *testing.T) {
	data := []int{1, 2, 3}
	last, err := lo.Last(data)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(last)

	last, err = lo.Last([]int{})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(last)
}

func TestNth(t *testing.T) {
	var (
		data = []int{0, 1, 2, 3, 4}
		nth  int
		err  error
	)

	nth, err = lo.Nth(data, 1)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(nth)

	nth, err = lo.Nth(data, -1)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(nth)
}

func TestSample(t *testing.T) {
	data := []string{"a", "b", "c"}

	sample := lo.Sample(data)
	fmt.Println(sample)
}

func TestSamples(t *testing.T) {
	data := []string{"a", "b", "c", "d", "e", "f", "g"}

	samples := lo.Samples(data, 2)
	fmt.Println(samples)
}
