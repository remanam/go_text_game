package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/remanam/go_text_game/game"
)

type GameState int

const (
	NEW GameState = iota
	STARTED
	FINISHED
)

type Game struct {
	Rooms     []game.Room
	Players   []game.Player
	GameState GameState
}

func (g GameState) String() string {
	switch g {
	case NEW:
		return "NEW"
	case STARTED:
		return "STARTED"
	case FINISHED:
		return "FINISHED"
	default:
		return "ERROR FUCK"
	}
}

func readRoomNames() {
	// Открываем файл для чтения
	file, err := os.Open("rooms.txt")
	if err != nil {
		fmt.Println("Ошибка при открытии файла:", err)
		return
	}
	defer file.Close()

	var roomNames []string

	// Создаем сканер для построчного чтения
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			roomNames = append(roomNames, line)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Ошибка при чтении файла:", err)
		return
	}

	// Выводим считанные названия
	fmt.Println("Названия комнат:")
	for _, name := range roomNames {
		fmt.Println(name)
	}
}

func initGame() {
	roomNames := readRoomNames()
	return Game{Rooms: roomNames, Players: []game.Player{}, GameState: NEW}
}

func main() {
	initGame()
}
