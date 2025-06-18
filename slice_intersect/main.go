package main

import (
	"fmt"
	"sort"
)

func intersection(a, b []int) []int {
	// Под капотом 3 вида сортировки, но по дефолту Quick sort
	sort.Ints(a) // N log N,  N = len(a)
	sort.Ints(b) // M log M, M = len(b)

	i := 0
	j := 0
	var result []int
	for i < len(a) && j < len(b) {
		if a[i] == b[i] {
			result = append(result, a[i])
			i++
			j++
		} else if a[i] < b[j] {
			i++
		} else {
			j++
		}
	}
	return result
}

func main() {

	a := []int{23, 3, 1, 2}
	b := []int{6, 2, 4, 23}
	// [2, 23]
	fmt.Printf("%v\n", intersection(a, b))
	a = []int{1, 1, 1}
	b = []int{1, 1, 1, 1}
	// [1, 1, 1]
	fmt.Printf("%v\n", intersection(a, b))
}
