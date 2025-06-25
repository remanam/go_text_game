package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/golang-jwt/jwt/v5"
)

// ЗАДАЧА 7. Работа с JWT в контексте
//Напиши две функции:
// AddJWTToContext(ctx context.Context, userID int) (context.Context, error)
// Эта функция должна -
// Создавать JWT-токен (используй github.com/golang-jwt/jwt/v5)
// Зашифровывать в него userID
// Возвращать новый context.Context, в который записан JWT-токен.

// ExtractUserIDFromContext(ctx context.Context) (int, error)
// Извлекать JWT-токен из контекста.
// Расшифровывать userID
// Вывести его на экран

// Дополнительное условие:
// Создай горутину, в которой будет использоваться ExtractUserIDFromContext.
//Покажи, что передача контекста работает и данные можно безопасно извлекать между горутинами.

const jwtSecret = "33rr3"

type jwtContextKey string

const jwtKey jwtContextKey = "jwtContextKey"

func AddJWTToContext(ctx context.Context, userID int) (context.Context, error) {
	claims := jwt.MapClaims{
		"userId": userID,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return context.Background(), err
	}

	context := context.WithValue(ctx, jwtKey, signedToken)

	return context, nil
}

func ExtractUserIDFromContext(ctx context.Context) (int, error) {
	token := ctx.Value(jwtKey)

	parts := strings.Split(token.(string), ".")
	if len(parts) != 3 {
		return 0, fmt.Errorf("invalid token format")
	}

	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return 0, err
	}
	var result map[string]int
	json.Unmarshal(payloadBytes, &result) // Переводим строку в map

	return result["userId"], nil
}

func main() {
	ctx := context.Background()
	userId := 1

	newCtx, err := AddJWTToContext(ctx, userId)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(newCtx)

	wg := sync.WaitGroup{}

	wg.Add(1)

	go func() {
		usId, err := ExtractUserIDFromContext(newCtx)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("Вытащенный айдишник = ", usId)
		wg.Done()
	}()

	wg.Wait()

}
