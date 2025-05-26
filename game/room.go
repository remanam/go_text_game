package game

type Room struct {
	Id         int
	Name       string
	RoomRoutes []string
	Commands   []string
}

type RoomActions struct {
	Name int // Основной идентификатор

}

// В какие комнаты можно попасть из текущей комнаты
type RoomRoutes struct {
	AvaiableRooms []Room
}
