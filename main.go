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

func readRoomNames() ([]string, error) {
	// Открываем файл для чтения
	file, err := os.Open("rooms.txt")
	if err != nil {
		fmt.Println("Ошибка при открытии файла:", err)
		return []string{}, err
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
		return []string{}, err
	}

	return roomNames, nil
}

var GameCore Game

func initGame() {
	roomNames, err := readRoomNames()
	if err != nil {
		fmt.Errorf("Ошибка при чтении названий комнат")
	}
	var rooms []game.Room
	for i := range roomNames {
		rooms = append(rooms, game.Room{
			Id:         i,
			Name:       roomNames[i],
			RoomRoutes: []string{},
			Commands:   []string{},
		})
	}
	GameCore = Game{Rooms: rooms, Players: []game.Player{}, GameState: NEW}
}

func main() {
	initGame()
}
