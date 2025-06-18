package main

import "fmt"

func mergeChannels(chans ...<-chan int) chan int {
	resultChan := make(chan int)

	go func() {
		for _, ch := range chans {
			for val := range ch {
				resultChan <- val
			}
		}
		close(resultChan)
	}()

	return resultChan
}

func main() {
	a := make(chan int)
	b := make(chan int)
	c := make(chan int)
	go func() {
		for _, val := range []int{1, 2, 3} {
			a <- val
		}
		close(a)
	}()
	go func() {
		for _, val := range []int{20, 10, 30} {
			b <- val
		}
		close(b)
	}()
	go func() {
		for _, val := range []int{300, 200, 100} {
			c <- val
		}
		close(c)
	}()

	for val := range mergeChannels(a, b, c) {
		fmt.Println(val)
	}
}
