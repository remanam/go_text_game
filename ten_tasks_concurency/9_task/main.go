package main

import (
	"fmt"
	"sync"
)

// Fan-In - это паттерн обмена сообщениями, используемый для создания воронки для работы среди исполнителей.
// Источником сообщений могут быть клиенты, а место назначения - сервер.
// Fan-In (объединение данных из нескольких каналов)

// ЗАДАЧА: Напишите функцию, которая объединяет три входных канала в ОДИН выходной.

// Объединяем разные каналы в один канал
func Merge(chs ...<-chan int) <-chan int {
	// chs - список каналов, которые объединим в один канал - result
	// Получается в чем суть - мы в каждый канал что-то записали в main
	// Затем здесь в горутине считали значения одного из каналов в out
	// Затем в main мы данные из этого канал выводим в Println

	wg := sync.WaitGroup{}
	wg.Add(3)

	// Создаем выходной канал
	out := make(chan int)

	for _, ch := range chs {
		go func(c <-chan int) {
			defer wg.Done()
			out <- <-ch
			fmt.Println("Прочитали с канала")
		}(ch)

	}
	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func main() {
	var wg sync.WaitGroup

	ch1 := make(chan int)
	ch2 := make(chan int)
	ch3 := make(chan int)

	fmt.Println("Создали 3 канала")

	wg.Add(3)
	go func() {
		ch1 <- 11
		fmt.Println("Записали в первый")
		close(ch1)
		wg.Done()
	}()
	go func() {
		ch2 <- 22
		close(ch2)
		fmt.Println("Записали во второй")
		wg.Done()
	}()
	go func() {
		ch3 <- 33
		close(ch3)
		fmt.Println("Записали в третий")
		wg.Done()
	}()
	merged := Merge(ch1, ch2, ch3)
	for val := range merged {
		fmt.Println(val)
	}

	wg.Wait()

}
