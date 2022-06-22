package lo_test

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/samber/lo"
	lop "github.com/samber/lo/parallel"
)

func TestMap(t *testing.T) {
	start := time.Now()
	m := lo.Map([]int64{1, 2, 3, 4}, func(x int64, _ int) string {
		return strconv.FormatInt(x, 10)
	})
	end := time.Now().Sub(start)
	fmt.Println(m, end)

	start2 := time.Now()
	m2 := lop.Map([]int64{1, 2, 3, 4, 5, 6, 7, 8, 9}, func(x int64, _ int) string {
		return strconv.FormatInt(x, 10)
	})
	end2 := time.Now().Sub(start2)
	fmt.Println(m2, end2)
}

func TestFlatMap(t *testing.T) {
	data := []int64{1, 2, 3, 4, 5, 6, 7, 8}

	fm := lo.FlatMap(data, func(x int64, _ int) []string {
		return []string{
			strconv.FormatInt(x, 10),
			strconv.FormatInt(x, 10),
			strconv.FormatInt(x, 10),
		}
	})

	fmt.Println(fm)
}

func TestFilterMap(t *testing.T) {
	data := []string{"cpu", "gpu", "mouse", "keyboard"}
	callback := func(x string, _ int) (string, bool) {
		if strings.HasSuffix(x, "pu") {
			return x, true
		}
		return "", false
	}

	matching := lo.FilterMap(data, callback)
	fmt.Println(matching)
}

func TestFilter(t *testing.T) {
	data := []int{1, 2, 3, 4}
	callback := func(x int, _ int) bool {
		return x%2 == 0
	}

	even := lo.Filter(data, callback)
	fmt.Println(even)
}

func TestContains(t *testing.T) {
	data := []int{1, 2, 3, 4, 5}
	present := lo.Contains(data, 5)

	fmt.Println(present)
}

func TestContainsBy(t *testing.T) {
	data := []int{1, 2, 3, 4, 5}
	callback := func(x int) bool {
		return x > 2 && x < 5
	}

	present := lo.ContainsBy(data, callback)
	fmt.Println(present)
}

func TestReduce(t *testing.T) {
	data := []int{1, 2, 3, 4}
	callback := func(agg int, item int, _ int) int {
		return agg + item*item
	}

	sum := lo.Reduce(data, callback, 1)
	fmt.Println(sum)
}

func TestForEach(t *testing.T) {
	data := []string{"hello", "world"}
	callback := func(x string, _ int) {
		fmt.Println(x)
	}
	lo.ForEach(data, callback)
	lop.ForEach(data, callback)
}

func TestTimes(t *testing.T) {
	data := 3
	callback := func(x int) string {
		return strconv.FormatInt(int64(x), 10)
	}
	r := lo.Times(data, callback)
	fmt.Printf("%T %v \n", r, r)
}

func TestUniq(t *testing.T) {
	data := []int{1, 2, 2, 1}
	uniq := lo.Uniq(data)
	fmt.Println(uniq)
}

func TestUniqBy(t *testing.T) {
	data := []int{0, 1, 2, 3, 4, 5, 6}
	callback := func(x int) int { return x % 3 }

	uniq := lo.UniqBy(data, callback)
	fmt.Println(uniq)
}

func TestGroupBy(t *testing.T) {
	data := []int{0, 1, 2, 3, 4, 5}
	callback := func(x int) int { return x % 3 }

	g := lo.GroupBy(data, callback)
	fmt.Println(g)

	gp := lop.GroupBy(data, callback)
	fmt.Println(gp)
}

func TestChunk(t *testing.T) {
	data := []int{0, 1, 2, 3, 4, 5}

	c1 := lo.Chunk(data, 2)
	fmt.Println(c1)

	data2 := []int{0, 1, 2, 3, 4, 5, 6}
	c2 := lo.Chunk(data2, 2)
	fmt.Println(c2)

	c3 := lo.Chunk([]int{}, 2)
	fmt.Println(c3)

	c4 := lo.Chunk([]string{"A"}, 2)
	fmt.Println(c4)
}

func TestPartitionBy(t *testing.T) {
	data := []int{-2, -1, 0, 1, 2, 3, 4, 5}
	callback := func(x int) string {
		if x < 0 {
			return "negative"
		} else if x%2 == 0 {
			return "even"
		}
		return "odd"
	}

	p := lo.PartitionBy(data, callback)
	fmt.Println(p)

	pp := lop.PartitionBy(data, callback)
	fmt.Println(pp)
}

func TestFlatten(t *testing.T) {
	data := [][]int{{0, 1}, {2, 3, 4, 5}}
	flat := lo.Flatten(data)
	fmt.Println(flat)
}

func TestShuffle(t *testing.T) {
	data := []int{1, 2, 3, 4, 5, 6, 7}
	sh := lo.Shuffle(data)
	fmt.Println(sh)
}

func TestReverse(t *testing.T) {
	data := []int{}
	for i := 0; i < 10; i++ {
		data = append(data, i)
	}

	t2 := time.Now()
	tmp := make([]int, len(data))
	length := len(data)
	for i := range tmp {
		tmp[i] = data[length-1-i]
	}
	e2 := time.Now().Sub(t2)
	fmt.Println("e2:", e2)

	lo.Reverse(data)
	t1 := time.Now()
	lo.Reverse(data)
	e1 := time.Now().Sub(t1)
	fmt.Println("e1:", e1)
}

type foo struct {
	bar string
}

func (f foo) Clone() foo {
	return foo{f.bar}
}
func TestFill(t *testing.T) {
	f := lo.Fill([]foo{{}, {}}, foo{"b"})
	fmt.Printf("%T %v", f, f)
}

func TestRepeat(t *testing.T) {
	slice := lo.Repeat(2, foo{"linonon"})
	fmt.Println(slice)
}

func TestRepeatBy(t *testing.T) {
	callback := func(x int) float64 {
		return math.Pow(float64(x), 2)
	}

	r := lo.RepeatBy(0, callback)
	fmt.Println(r)

	r2 := lo.RepeatBy(5, callback)
	fmt.Println(r2)
}

type Character struct {
	dir  string
	code int
}

func TestKeyBy(t *testing.T) {
	// 1
	data := []string{"a", "aa", "aaa"}
	callback := func(str string) int {
		return len(str)
	}

	m := lo.KeyBy(data, callback)
	fmt.Println(m)

	// 2
	characters := []Character{
		{dir: "left", code: 97},
		{dir: "right", code: 100},
	}

	callback2 := func(char Character) string {
		return string(rune(char.code))
	}

	result := lo.KeyBy(characters, callback2)
	fmt.Println(result)
}

func TestDrop(t *testing.T) {
	data := []int{0, 1, 2, 3, 4, 5, 6}
	l := lo.Drop(data, 2)
	fmt.Println(l)
}

func TestDropRight(t *testing.T) {
	data := []int{0, 1, 2, 3, 4, 5, 6}

	dr := lo.DropRight(data, 2)
	fmt.Println(dr)
}

func TestDropWhile(t *testing.T) {
	data := []string{"a", "aa", "aaa", "aa", "aa"}
	callback := func(str string) bool { return len(str) <= 2 }

	dw := lo.DropWhile(data, callback)
	fmt.Println(dw)
}

func TestDropRightWhile(t *testing.T) {
	data := []string{"a", "aa", "aaa", "aa", "aa"}
	callback := func(str string) bool { return len(str) <= 2 }

	drw := lo.DropRightWhile(data, callback)
	fmt.Println(drw)
}

func TestReject(t *testing.T) {
	data := []int{1, 2, 3, 4}
	callback := func(x int, _ int) bool { return x%2 == 0 }

	odd := lo.Reject(data, callback)
	fmt.Println(odd)
}

func TestCounters(t *testing.T) {
	data := []int{1, 2, 3, 4, 1, 1, 2}
	count := lo.Count(data, 1)
	fmt.Println(count)
}

func TestCountBy(t *testing.T) {
	data := []int{1, 2, 3, 4, 5}
	callback := func(x int) bool { return x <= 3 }

	count := lo.CountBy(data, callback)
	fmt.Println(count)
}

func TestSubset(t *testing.T) {
	data := []int{0, 1, 2, 3, 4, 5}

	sub := lo.Subset(data, 2, 3)
	fmt.Println(sub)

	sub2 := lo.Subset(data, -4, 3)
	fmt.Println(sub2)

	sub3 := lo.Subset(data, -2, math.MaxUint)
	fmt.Println(sub3)
}

func TestReplace(t *testing.T) {
	data := []int{0, 1, 0, 1, 2, 3, 0}

	var s []int
	s = lo.Replace(data, 0, 42, 1)
	fmt.Println(s)

	s = lo.Replace(data, -1, 42, 1)
	fmt.Println(s)

	s = lo.Replace(data, 0, 42, 2)
	fmt.Println(s)

	s = lo.Replace(data, 0, 42, -1)
	fmt.Println(s)
}

func TestReplaceAll(t *testing.T) {
	data := []int{0, 1, 0, 1, 2, 3, 0}
	var s []int

	s = lo.ReplaceAll(data, 0, 42)
	fmt.Println(s)

	s = lo.ReplaceAll(data, -1, 42)
	fmt.Println(s)
}

func TestKeys(t *testing.T) {
	data := map[string]int{"foo": 1, "bar": 2}

	keys := lo.Keys(data)
	fmt.Println(keys)
}

func TestValues(t *testing.T) {
	data := map[string]int{"foo": 1, "bar": 2}

	values := lo.Values(data)
	fmt.Println(values)
}

func TestPickBy(t *testing.T) {
	data := map[string]int{"foo": 1, "bar": 2, "baz": 3}
	callback := func(key string, value int) bool { return value%2 == 1 }

	m := lo.PickBy(data, callback)
	fmt.Println(m)
}

func TestPickByKeys(t *testing.T) {
	data := map[string]int{"foo": 1, "bar": 2, "baz": 3}
	keys := []string{"foo", "bar"}
	m := lo.PickByKeys(data, keys)
	fmt.Println(m)
}

func TestPickByValues(t *testing.T) {
	data := map[string]int{"foo": 1, "bar": 2, "baz": 3, "bal": 3}
	keys := []int{1, 3}
	m := lo.PickByValues(data, keys)
	fmt.Println(m)
}
