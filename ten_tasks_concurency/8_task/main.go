package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

// ЗАДАЧА:
// Напишите функцию FetchURLs(urls []string) map[string]string, которая:

// Принимает слайс URL-адресов.
// Конкурентно делает HTTP-запросы к каждому URL.
// Собирает результаты (код ответа и часть тела) в map[string]string, где:
// ключ — URL
// значение — содержимое ответа (ограниченное, например, 100 символами)
// Использует sync.WaitGroup и sync.Mutex для защиты записи в map.
// В случае ошибки записывает "error" как значение.

// *Скорее всего, ты сделал без таймаутов на запрос и контекстов. Надо это исправить
// DONE

// TODO
//Средний уровень:
// 2. Переписать код, чтобы использовать контекст отмены
// 3. При первой ошибке отменять все остальные запросы

// TODO
//Продвинутый уровень(опционально):
//4. Ограничить общее количество параллельных запросов (например, максимум 10 одновременно через семафор — chan struct{}).
// А лучше использовать для этого паттерн worker pool Go by Example: Worker Pools

func FetchURLs(urls []string, wg *sync.WaitGroup, ctx context.Context) map[string]string {
	client := http.Client{Timeout: 10 * time.Second}

	result := make(map[string]string)

	mutex := sync.Mutex{}

	for i := range urls {
		go func() {

			req, err := http.NewRequestWithContext(ctx, http.MethodGet, urls[i], nil)
			if err != nil {
				fmt.Printf("Ошибка создания запроса для url %s\n", urls[i])

			}

			resp, err := client.Do(req)
			if err != nil {
				fmt.Println(err.Error())
			}
			// Неверный статус ответа
			if resp.StatusCode != http.StatusOK {
				fmt.Println("Неверный статус для запроса: ", urls[i])
			}
			limited := io.LimitReader(resp.Body, 30)
			body, err := io.ReadAll(limited)
			// Если не смогли прочитать тело
			if err != nil {
				fmt.Println(err.Error())
			}
			defer resp.Body.Close()
			mutex.Lock()
			result[urls[i]] = string(body)
			mutex.Unlock()

			wg.Done()

		}()
	}

	return result
}

func main() {
	urls := []string{
		"https://playground.epizy.com/",
		"https://practice-automation.com/",
		"https://automation-playground.com/",
		"https://demoqa.com/",
		"https://reqres.in/",
		"https://rickandmortyapi.com/",
	}
	wg := sync.WaitGroup{}
	wg.Add(len(urls))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	result := FetchURLs(urls, &wg, ctx)
	wg.Wait()
	for key, value := range result {
		fmt.Printf("url %s = %s \n", key, value)
	}

}
