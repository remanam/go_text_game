package game

type Player struct {
	Id     int
	Name   string
	RoomId int
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
