package models

type PlayersModel struct {
	Players []*Player `json:"players"`
}

type Player struct {
	PlayerId int    `json:"id"`
	Name     string `json:"name"`
	TeamId   int    `json:"teamId"`
	Role     string `json:"role"`
}

type PlayerStats struct {
	PlayerId int `json:"playerId"`
	Votes    int `json:"votes"`
}
