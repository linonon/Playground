package lo_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/samber/lo"
)

func TestOmitBy(t *testing.T) {
	data := map[string]int{"foo": 1, "bar": 2, "baz": 3, "bal": 3}
	callback := func(k string, v int) bool {
		return k == "foo" || v == 3
	}

	m := lo.OmitBy(data, callback)
	fmt.Println(m)
}

func TestOmitByKeys(t *testing.T) {
	data := map[string]int{"foo": 1, "bar": 2, "baz": 3, "bal": 3}
	keys := []string{"foo", "bar"}

	m := lo.OmitByKeys(data, keys)
	fmt.Println(m)
}

func TestOmitByValues(t *testing.T) {
	data := map[string]int{"foo": 1, "bar": 2, "baz": 3, "bal": 3}
	values := []int{1, 2}

	m := lo.OmitByValues(data, values)
	fmt.Println(m)
}

func TestEntries(t *testing.T) {
	data := map[string]int{"foo": 1, "bar": 2, "baz": 3, "bal": 3}

	e := lo.Entries(data)
	fmt.Printf("%+v", e)
}

func TestFromEntries(t *testing.T) {
	data := map[string]int{"foo": 1, "bar": 2, "baz": 3, "bal": 3}

	e := lo.Entries(data)

	data2 := lo.FromEntries(e)
	fmt.Println(data2)
}

func TestInvert(t *testing.T) {
	data1 := map[string]int{"a": 1, "b": 2, "c": 3}
	m1 := lo.Invert(data1)
	fmt.Println(m1)

	data2 := map[string]int{"a": 1, "b": 2, "c": 1, "d": 1, "e": 1, "f": 1, "g": 1}
	m2 := lo.Invert(data2)
	fmt.Println(m2)
}

func TestAssign(t *testing.T) {
	m1 := map[string]int{"a": 1, "b": 2}
	m2 := map[string]int{"b": 22, "c": 3}
	m3 := map[string]int{"b": 33}

	mm := lo.Assign(m1, m2, m3)
	fmt.Println(mm)
}

func TestMapKeys(t *testing.T) {
	data := map[int]int{1: 1, 2: 2, 3: 3, 4: 4, 5: 5}
	callback := func(_ int, v int) string {
		return strconv.FormatInt(int64(v), 10)
	}

	m := lo.MapKeys(data, callback)
	fmt.Printf("%+v", m)
}

func TestMapValues(t *testing.T) {
	data := map[int]int{1: 1, 2: 2, 3: 3, 4: 4, 5: 5}
	callback := func(k int, _ int) string {
		return strconv.FormatInt(int64(k), 10)
	}

	m := lo.MapValues(data, callback)
	fmt.Println(m)
}
