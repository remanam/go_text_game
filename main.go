package main

import (
	"fmt"
	"slices"
	"strings"
)

type GameState int

const (
	NEW GameState = iota
	STARTED
	FINISHED
)

type Game struct {
	Locations []Location
	Player    Player
	Step      int
	GameState GameState
	Quests    []Quest
}

func (g *Game) setCurrentRoom(roomName string) {
	for i := range g.Locations {
		if g.Locations[i].Name == roomName {
			g.Player.CurrentLocation = g.Locations[i]
		}
	}
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

// -------------------------------------- Комнаты начало
type Location struct {
	Name              string
	LocationRoutes    []string // куда можно пойти с этой комнаты
	LookupDescription string   //Описание ответа после команды "осмотреться"
	WelcomeMessage    string   //Описание ответа после команда "идти"
	PassRequirement   *Requirement
	Places            []Place
}

type Place struct {
	Name       string      // Например "На столе", "На стуле"
	PlaceItems []PlaceItem // Например "рюкзак, ключи"
	IsEmpty    bool        // Буду использовать этот флаг, когда PlaceItems будет пустой массив

}

type PlaceItem struct {
	Name     string
	Avaiable bool
}

func (g *Game) GetLocationByName(Name string) (Location, error) {

	for i := 0; i < len(g.Locations); i++ {
		if g.Locations[i].Name == Name {
			return g.Locations[i], nil
		}
	}

	return Location{}, fmt.Errorf("location %q not found", Name)

}

func (g *Game) AreQuestsDone() bool {
	for i := range g.Quests {
		if !g.Quests[i].Done {
			return false
		}
	}
	return true
}

//--------------------------------------- Комнаты конец

// игровые предметы
type Item map[string]int

type Requirement struct {
	ItemName         string
	LockName         string
	LocationToBlock  Location //Локация, в которую из этой комнаты нельзя пойти без requirement
	Satisfied        bool
	ForbiddenMessage string // Текст, если пытаемся пройти в комнату, когда Satisfied = false
}

type Quest struct {
	Name string
	Done bool
}

// --------------------------------------- Игрок начало
type Player struct {
	Id              int
	Name            string
	CurrentLocation Location
	Inventory       Item
}

func (p *Player) CheckItemInPlayerStock(itemName string) bool {
	v, ok := p.Inventory[itemName]
	if !ok {
		//У игрока нет нужной вещи
		return false
	}
	// Если количество <=0
	if v <= 0 {
		return false
	}
	// Если предмет у игрока есть
	return true

}

// --------------------------------------------- Игрок Конец

var UNKNOWN_COMMAND string = "неизвестная команда"

func lookCommand() string {
	currLocation := GameExample.Player.CurrentLocation

	var routes string
	for i := 0; i < len(currLocation.LocationRoutes); i++ {
		if i > 0 {
			routes += ", "
		}
		routes += currLocation.LocationRoutes[i]
	}

	// Здесь мы помечаем какие Place у нас остались БЕЗ предметов
	// Ниже для расчет places будем использовать
	for i := range currLocation.Places {
		if len(currLocation.Places[i].PlaceItems) == 0 {
			currLocation.Places[i].IsEmpty = true
		}
	}

	isLocationEmpty := true
	var places string
	// if currLocation.Name == "комната" {

	for i := range currLocation.Places {
		if !currLocation.Places[i].IsEmpty {
			isLocationEmpty = false
		}
	}

	for i := range len(currLocation.Places) {
		if isLocationEmpty && currLocation.Name == "комната" {
			places = "пустая комната"
			break
		}
		if currLocation.Places[i].IsEmpty {
			// Если в этом Place нет предметов, то НЕ выводим
			continue
		}

		places += currLocation.Places[i].Name
		var items []string
		for j := range len(currLocation.Places[i].PlaceItems) {
			items = append(items, currLocation.Places[i].PlaceItems[j].Name)
		}
		places += strings.Join(items, ", ")
		//К последнему place запятую НЕ ставим
		if i < len(currLocation.Places)-1 &&
			currLocation.Places[len(currLocation.Places)-1].IsEmpty == false {
			places += ", "
		}
	}

	var quests string
	isFirst := true
	if !GameExample.AreQuestsDone() {
		for i := range GameExample.Quests {
			if !GameExample.Quests[i].Done {
				if isFirst {
					quests += "надо " + GameExample.Quests[i].Name
					isFirst = false
				} else {
					quests += " и " + GameExample.Quests[i].Name
				}
			}

		}
	}

	if currLocation.Name != "кухня" {
		// Квесты НЕ показываем
		return currLocation.LookupDescription + places + "." + " можно пройти - " + routes
	} else {
		return currLocation.LookupDescription + places + ", " + quests + "." + " можно пройти - " + routes
	}
}

func goCommand(locationName string) string {
	if !slices.Contains(GameExample.Player.CurrentLocation.LocationRoutes, locationName) {
		return "нет пути в " + locationName
	}
	currLocation := GameExample.Player.CurrentLocation
	locationToGo, err := GameExample.GetLocationByName(locationName)
	if err != nil {
		return "Ошибки в получении локации"
	}

	if currLocation.PassRequirement != nil {
		if currLocation.PassRequirement.LocationToBlock.Name == locationToGo.Name && !currLocation.PassRequirement.Satisfied {
			return currLocation.PassRequirement.ForbiddenMessage
		}
	}

	var newRoutes string
	GameExample.setCurrentRoom(locationToGo.Name)
	for i := 0; i < len(locationToGo.LocationRoutes); i++ {
		if i > 0 {
			newRoutes += ", "
		}
		newRoutes += locationToGo.LocationRoutes[i]
	}
	fmt.Println("-----")

	if GameExample.Player.CurrentLocation.PassRequirement == nil {
		GameExample.setCurrentRoom(locationToGo.Name)
		return GameExample.Player.CurrentLocation.WelcomeMessage + " можно пройти - " + newRoutes
	}

	GameExample.setCurrentRoom(locationToGo.Name)
	return GameExample.Player.CurrentLocation.WelcomeMessage + " можно пройти - " + newRoutes

	// ничего интересного. можно пройти - кухня, комната, улица

}

func handleCommand(command string) string {
	//
	parts := strings.Split(command, " ")
	fmt.Println(parts)
	switch parts[0] {
	case "осмотреться":
		return lookCommand()
	case "идти":
		return goCommand(parts[1])
	case "надеть":
		return pickupCommand(parts[1])
	case "взять":
		return pickupCommand(parts[1])
	case "применить":
		return applyCommand(parts[1], parts[2])
	}
	return "неизвестная команда"
}

func applyCommand(itemName string, lockName string) string {
	player := GameExample.Player
	if player.CheckItemInPlayerStock(itemName) {
		if player.CurrentLocation.PassRequirement != nil {
			if lockName == player.CurrentLocation.PassRequirement.LockName {
				player.CurrentLocation.PassRequirement.Satisfied = true
				return "дверь открыта"
			} else {
				return "не к чему применить"
			}
		}
	} else {
		return "нет предмета в инвентаре - " + itemName
	}
	return "Ошибка в команде apply"
}

func pickupCommand(itemName string) string {
	isItemInLocation := false
	for i := range GameExample.Player.CurrentLocation.Places {
		for j := range GameExample.Player.CurrentLocation.Places[i].PlaceItems {
			// Если такой предмет существует на локации

			if GameExample.Player.CurrentLocation.Places[i].PlaceItems[j].Name == itemName &&
				GameExample.Player.CurrentLocation.Places[i].PlaceItems[j].Avaiable {

				if !GameExample.Player.CheckItemInPlayerStock("рюкзак") && itemName != "рюкзак" {
					return "некуда класть"
				}

				GameExample.Player.Inventory[itemName] = 1
				isItemInLocation = true

				GameExample.Player.CurrentLocation.Places[i].PlaceItems = append(
					GameExample.Player.CurrentLocation.Places[i].PlaceItems[:j],
					GameExample.Player.CurrentLocation.Places[i].PlaceItems[j+1:]...)

				if itemName == "рюкзак" {
					for i := range GameExample.Quests {
						if GameExample.Quests[i].Name == "собрать рюкзак" {
							GameExample.Quests[i].Done = true
						}
					}
					return "вы надели: " + itemName
				} else if GameExample.Player.CheckItemInPlayerStock("рюкзак") {
					// Поднимать любые предметы, кроме рюкзака можно после надевания рюкзака
					GameExample.Player.Inventory[itemName] = 1
					return "предмет добавлен в инвентарь: " + itemName
				}
			}
		}
	}
	if !isItemInLocation {
		return "нет такого"
	}
	return "ошибка в команде pickup"
}

var GameExample Game

func initGame() {
	//ЛОКАЦИИ "кухня", "коридор", "комната", "улица",  "домой"}

	kitchen := Location{
		Name:              "кухня",
		LocationRoutes:    []string{"коридор"},
		LookupDescription: "ты находишься на кухне, ",
		WelcomeMessage:    "кухня, ничего интересного.", //Это текст, который показывается после перехода в эту комнату
		PassRequirement:   nil,
		Places:            []Place{{Name: "на столе: ", PlaceItems: []PlaceItem{{Name: "чай", Avaiable: true}}}}}
	street := Location{
		Name:              "улица",
		LocationRoutes:    []string{"домой"},
		LookupDescription: "",
		WelcomeMessage:    "на улице весна.",
		PassRequirement:   nil,
		Places:            []Place{},
	}
	corridor := Location{
		Name:              "коридор",
		LocationRoutes:    []string{"кухня", "комната", "улица"},
		LookupDescription: "",
		WelcomeMessage:    "ничего интересного.",
		PassRequirement:   &Requirement{ItemName: "ключи", LockName: "дверь", LocationToBlock: street, Satisfied: false, ForbiddenMessage: "дверь закрыта"},
		Places:            []Place{},
	}

	room := Location{
		Name:              "комната",
		LocationRoutes:    []string{"коридор"},
		LookupDescription: "",
		WelcomeMessage:    "ты в своей комнате.",
		PassRequirement:   nil,
		Places: []Place{
			{Name: "на столе: ", PlaceItems: []PlaceItem{{Name: "ключи", Avaiable: true}, {Name: "конспекты", Avaiable: true}}},
			{Name: "на стуле: ", PlaceItems: []PlaceItem{{Name: "рюкзак", Avaiable: true}}},
		},
	}

	Locations := []Location{kitchen, corridor, room, street}

	quests := []Quest{
		{Name: "собрать рюкзак", Done: false}, {Name: "идти в универ", Done: false},
	}
	player := Player{
		Id:              0,
		Name:            "chlenVPalto",
		CurrentLocation: kitchen,
		Inventory:       make(Item),
	}

	GameExample = Game{Locations: Locations, Player: player, GameState: NEW, Quests: quests}
	GameExample.Player.CurrentLocation = kitchen
}
