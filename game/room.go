package game

type roomStatus string

const (
	StatusOnHold  = roomStatus("STATUS_ON_HOLD")
	StatusRunning = roomStatus("STATUS_RUNNING")
)

type Room interface {
	Enter(player Player)
	Exit(player Player)
}

type tiles struct {
	horizontal int
	vertical   int
}

type room struct {
	ID           uint64
	playersLimit int
	status       roomStatus
	owner        Player
	players      []Player
	tiles        tiles
}

func NewRoom() Room {
	return &room{}
}

func (r room) playersLen() int {
	return len(r.players) + 1
}

func (r room) Enter(player Player) {
	if r.owner.GetID() != player.GetID() && r.playersLen() < int(r.playersLimit) {
		r.players = append(r.players, player)
	}
}

func (r room) Exit(player Player) {

}
