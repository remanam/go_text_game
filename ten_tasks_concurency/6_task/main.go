package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

//ЗАДАЧА: Батчевая обработка
//Реализуй функцию StartBatchProcessor(ctx context.Context, input <-chan int), которая:
// 1) Собирает числа из канала input в батчи по максимум 5 элементов.
// 2) Если в течение 2 секунд батч не собран — обрабатывает то, что есть.
// 3) Обработка батча — это просто fmt.Println("Processed batch:", batch).
// 4) Выход из функции должен происходить при отмене контекста (ctx.Done()).
// Дополнительно Отмена должна происходить либо через context.WithTimeout, либо вручную через cancel() — попробовать оба варианта

func StartBatchProcessor(ctx context.Context, input chan int, wg *sync.WaitGroup) {
	// Поскольку есть у канала буфер или нет мы не знаем, то пишу здесь
	// буфер будет 5 элементов
	const batchSize = 5
	const timerDuration = 2 * time.Second
	batch := make([]int, 0, batchSize)
	timer := time.NewTimer(timerDuration)

	for {
		select {
		case <-ctx.Done():
			wg.Done()
			fmt.Println("Горутина закончила выполнение")
			return
		case <-timer.C:
			fmt.Println("За 2 секунды успело обработаться ", len(batch), " элементов.")

			batch = make([]int, 0, batchSize)
			fmt.Println("Заново создали BATCH - С ТАЙМЕРОМ")
		case val := <-input:
			if len(batch) < 5 {
				batch = append(batch, val)
				fmt.Println("Processed batch: ", val)
			} else if len(batch) == 5 {
				batch = make([]int, 0, batchSize)
				fmt.Println("Заново создали BATCH - БЕЗ ТАЙМЕРА")
				timer.Reset(2 * time.Second)
			}

		}

	}
}

func main() {
	// инициализация канала
	/* создание контекста  */
	wg := sync.WaitGroup{}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	input := make(chan int)
	wg.Add(1)

	go func() {
		for range 20 {
			random := rand.Intn(100)
			input <- random
			time.Sleep(150 * time.Millisecond)
		}
	}()

	go StartBatchProcessor(ctx, input, &wg)

	// сбор данных
	wg.Wait()

	fmt.Println("Main: processing stopped")
}
