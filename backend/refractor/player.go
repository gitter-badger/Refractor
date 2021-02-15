package refractor

type Player struct {
	PlayerID      int64    `json:"id"`
	PlayFabID     string   `json:"playFabId"`
	LastSeen      int64    `json:"lastSeen"`
	CurrentName   string   `json:"currentName"`
	PreviousNames []string `json:"previousNames,omitempty"`
}

type PlayerRepository interface {
	Create(player *Player) error
	FindByID(id int64) (*Player, error)
	FindByPlayFabID(playFabID string) (*Player, error)
	Exists(args FindArgs) (bool, error)
	UpdateName(player *Player, currentName string) error
	Update(id int64, args UpdateArgs) (*Player, error)
}

type PlayerService interface {
	CreatePlayer(newPlayer *Player) (*Player, *ServiceResponse)
}