package main

import "fmt"

func main() {
	naturals := make(chan int)
	squares := make(chan int)

	go func() {
		for x := 0; x < 10; x++ {
			naturals <- x
		}
		close(naturals)
	}()

	go func() {

		for x := range naturals {
			go func(x int) {
				squares <- x * x
			}(x)

		}
		close(squares)
	}()

	for x := range squares {
		fmt.Println(x)
	}
}
