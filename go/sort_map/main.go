package main

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

// FrequencyList is a list with frequency, the order is decided by `sortOrder`
//
// Example:
//
//	fl := NewFrequencyList(0, 2, -1)
//	fl.Append(NewFrequencyItem(0, 1))
//	fl.Append(NewFrequencyItem(1, 2))
//	fl.Append(NewFrequencyItem(2, 3))
//	fl.Append(NewFrequencyItem(3, 1))
//	fl.Sort()
//	assert.Equal(fl.ToIntSlice(), []int{2, 1, 0, 3})
type FrequencyList struct {
	list      []frequencyItem
	sortOrder int
}

// NewFrequencyList create a ListSlice with [start, start+1, ..., end-1, end] elements,
// if order == -1, the list will be sorted in descending order, otherwise ascending order
//
// Example:
//
//	assert.Equal(NewFrequencyList(0, 2,  1), []int{0, 1, 2})
//	assert.Equal(NewFrequencyList(0, 2, -1), []int{2, 1, 0})
func NewFrequencyList(start, end int, order int) FrequencyList {
	s := make([]frequencyItem, end-start+1)
	for i := start; i <= end; i++ {
		s[i] = NewFrequencyItem(i)
	}

	return FrequencyList{
		list:      s,
		sortOrder: order,
	}
}

// Append item to the list, item can be created by NewFrequencyItem(k, v)
func (fl *FrequencyList) Append(item ...frequencyItem) {
	fl.list = append(fl.list, item...)
}

// Remove item from the list
func (fl *FrequencyList) Remove(value int) {
	for i := range fl.list {
		if fl.list[i].Value == value {
			fl.list = append(fl.list[:i], fl.list[i+1:]...)
			break
		}
	}
}

// GetSortOrder return the sort order
func (fl *FrequencyList) GetSortOrder() int {
	return fl.sortOrder
}

func (fl *FrequencyList) SetSortOrder(order int) {
	fl.sortOrder = order
}

// Update the frequent of item with value
func (fl *FrequencyList) Update(value int) {
	for i := range fl.list {
		if fl.list[i].Value == value {
			fl.list[i].Freqent++
			fl.list[i].updateTimestamp = time.Now().UnixMicro()
			break
		}
	}
	fl.Sort()
}

// GetTopItem return the top k items
func (fl *FrequencyList) GetTopItem(k int) []frequencyItem {
	if k > fl.Len() {
		return fl.list
	}

	return append([]frequencyItem(nil), fl.list[:k]...)
}

// GetTailItem return the tail k items
func (fl *FrequencyList) GetTailItem(k int) []frequencyItem {
	if k > fl.Len() {
		return fl.list
	}

	return append([]frequencyItem(nil), fl.list[fl.Len()-k:]...)
}

// GetTopValue return the top k values
func (fl *FrequencyList) GetTopValue(k int) []int {
	if k > fl.Len() {
		return fl.ToIntSlice()
	}

	return fl.ToIntSlice()[:k]
}

// GetTailValue return the tail k values
func (fl *FrequencyList) GetTailValue(k int) []int {
	if k > fl.Len() {
		return fl.ToIntSlice()
	}

	return fl.ToIntSlice()[fl.Len()-k:]
}

func (fl *FrequencyList) ToIntSlice() []int {
	res := make([]int, fl.Len())
	for i := range fl.list {
		res[i] = fl.list[i].Value
	}
	return res
}

// Len for sort.Interface
func (fl *FrequencyList) Len() int {
	return len(fl.list)
}

// Less for sort.Interface
func (fl *FrequencyList) Less(i, j int) bool {
	l := fl.list

	if fl.sortOrder < 0 {
		if l[i].Freqent == l[j].Freqent {
			return l[i].updateTimestamp > l[j].updateTimestamp
		} else {
			return l[i].Freqent > l[j].Freqent
		}
	} else {
		if l[i].Freqent == l[j].Freqent {
			return l[i].updateTimestamp < l[j].updateTimestamp
		} else {
			return l[i].Freqent < l[j].Freqent
		}
	}
}

// Swap for sort.Interface
func (fl *FrequencyList) Swap(i, j int) {
	l := fl.list
	l[i], l[j] = l[j], l[i]
}

// Shuffle the list
func (fl *FrequencyList) Shuffle() {
	l := fl.list
	rand.Shuffle(len(l), func(i, j int) { l[i], l[j] = l[j], l[i] })
}

// Sort the list, the order is decided by `fl.sortOrder`
func (fl *FrequencyList) Sort() {
	sort.Sort(fl)
}

type frequencyItem struct {
	Value   int
	Freqent int

	updateTimestamp int64
}

// NewFrequencyItem create a FrequencyItem with value
func NewFrequencyItem(value int) frequencyItem {
	return frequencyItem{
		Value:           value,
		Freqent:         0,
		updateTimestamp: time.Now().UnixMicro(),
	}
}

func main() {
	a := map[int]int{
		1:  1,
		2:  2,
		3:  3,
		4:  1,
		5:  1,
		6:  1,
		7:  1,
		8:  2,
		9:  3,
		10: 3,
		11: 4,
	}

	s := FrequencyList{}

	for k, v := range a {
		s.Append(frequencyItem{k, v, 0})
	}

	s.Sort()
	fmt.Println(s)

	// b
	b := NewFrequencyList(0, 20, -1)
	fmt.Println(b)

	b.Shuffle()
	fmt.Println(b)

	b.Update(1)
	b.Update(11)
	fmt.Println(b)

	fmt.Println("top 3:", b.GetTopItem(3))
	fmt.Println("tail 3:", b.GetTailItem(3))
}
