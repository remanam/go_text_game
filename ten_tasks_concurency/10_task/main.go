package main

import (
	"fmt"
	"math/rand"
	"sync"
)

// Fan-Out (разделение работы между несколькими воркерами)
// Задача: Создайте функцию, которая принимает канал с задачами и распределяет их между N горутинами.

// Разделить канал на n каналов,
// которые получают сообщения в циклическом порядке.
func Split(input <-chan int, n int, wg *sync.WaitGroup) []chan int {

	resultChans := make([]chan int, n)
	for i := range n {
		resultChans[i] = make(chan int)
	}
	go func() {
		for val := range input {
			randomIndex := rand.Intn(n)
			resultChans[randomIndex] <- val
			wg.Done()
		}

		// Когда все значения вытащили, то закрываем выходные каналы
		for _, val := range resultChans {
			close(val)
		}

	}()

	return resultChans
}

func main() {

	wg := sync.WaitGroup{}

	input := make(chan int)
	goroutinesCount := 4

	// 20 раз записываем, значит WG нужен 20
	wg.Add(20)
	go func() {

		for i := range 20 {

			input <- i * 10
		}
		// После того как записали все что хотели закрываем канал
		// ВОПРОС - почему без закрытия канала тоже все хорошо работает ?
		close(input)
	}()

	chans := Split(input, goroutinesCount, &wg)
	for i, ch := range chans {

		go func(i int, ch <-chan int) {
			for val := range ch {
				fmt.Println("Горутина - ", i, " обработала input - ", val)
			}
		}(i, ch)

	}

	wg.Wait()

}
