package lo_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/samber/lo"
)

func TestMust(t *testing.T) {
	{
		val := lo.Must(time.Parse("2006-01-02", "2002-01-14"))
		fmt.Println(val)
	}
	{
		val := lo.Must(time.Parse("2006-01-02", "bad-value"))
		fmt.Println(val)
	}

	// func example0() (error)
	// func example1() (int, error)
	// func example2() (int, string, error)
	// func example3() (int, string, time.Date, error)
	// func example4() (int, string, time.Date, bool, error)
	// func example5() (int, string, time.Date, bool, float64, error)
	// func example6() (int, string, time.Date, bool, float64, byte, error)

	// lo.Must0(example0())
	// val1 := lo.Must1(example1())    // alias to Must
	// val1, val2 := lo.Must2(example2())
	// val1, val2, val3 := lo.Must3(example3())
	// val1, val2, val3, val4 := lo.Must4(example4())
	// val1, val2, val3, val4, val5 := lo.Must5(example5())
	// val1, val2, val3, val4, val5, val6 := lo.Must6(example6())
}

func TestTry(t *testing.T) {
	{
		ok := lo.Try(func() error {
			panic("should not be called")
		})
		fmt.Println(ok)
	}
	{
		ok := lo.Try(func() error {
			return nil
		})
		fmt.Println(ok)
	}
	{
		ok := lo.Try(func() error {
			return fmt.Errorf("error")
		})
		fmt.Println(ok)
	}
}

func TestCatch(t *testing.T) {
	caught := false

	lo.TryCatch(func() error {
		panic("error")
	}, func() {
		caught = true
	})

	fmt.Println(caught)
}

func TestTryWithErrorValue(t *testing.T) {
	err, ok := lo.TryWithErrorValue(func() error {
		panic("error")
	})

	fmt.Println(err, ok)
}

func TestTryCatchWithErrorValue(t *testing.T) {
	caught := false

	lo.TryCatchWithErrorValue(func() error {
		panic("error")
	}, func(val any) {
		caught = val == "error"
	})

	fmt.Println(caught)
}
