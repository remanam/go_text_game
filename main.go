package main

import "slices"

type GameState int

const (
	NEW GameState = iota
	STARTED
	FINISHED
)

type Game struct {
	Rooms       []Room
	CurrentRoom Room
	Players     Player
	Step        int
	GameState   GameState
}

func (g *Game) setCurrentRoom(roomName string) {
	g.CurrentRoom = roomName
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

var GameExample Game

// -------------------------------------- Комнаты начало
type Room struct {
	Id         int
	Name       string
	RoomRoutes []string // куда можно пойти с этой комнаты
	Commands   []string // какие есть действия. Взять рюкзак типа
}

func (r *Room) GetAvaiableRooms(roomName string) []string {
	return r.RoomRoutes
}

func (r *Room) GetAvaiableCommands(roomName string) []string {
	return r.Commands
}

//--------------------------------------- Комнаты конец

// --------------------------------------- Игрок начало
type Player struct {
	Id              int
	Name            string
	CurrentRoomId   int
	AvaiableActions int
	inventory       []string
}

type PlayerActions interface {
	Move(RoomId int) error
	Action(ActionId int) error
}

func (p *Player) Move(RoomId int) error {
	return nil
}

func (p *Player) Action(ActionId int) error {
	return nil
}

// --------------------------------------------- Игрок Конец

var UNKNOWN_COMMAND string = "неизвестная команда"

func handleCommand(command string) string {
	if slices.Contains(GameExample.CurrentRoom.Commands, command) == false {
		return ""
	}
	if command == "осмотреться" {
		return "ты находишься на кухне, на столе: чай, надо собрать рюкзак и идти в универ. можно пройти - коридор"
	} else if command == "идти коридор" {
		return ""
	} else {
		return UNKNOWN_COMMAND

	}
}

func initGame() {
	roomNames := []string{"кухня", "универ", "коридор", "улица", "комната", "дом"}

	var rooms []Room
	for i := range roomNames {
		rooms = append(rooms, Room{
			Id:         i,
			Name:       roomNames[i],
			RoomRoutes: []string{},
			Commands:   []string{},
		})
	}
	GameExample = Game{Rooms: rooms, Players: Player{}, GameState: NEW}
}

func main() {
	initGame()
}
