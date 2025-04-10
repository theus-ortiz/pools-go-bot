package pools

type Position struct {
	ID      string `json:"id"`
	Address string `json:"address"`
	Network string `json:"network"`
}

type UserPools struct {
	Owner     string     `json:"owner"`
	Positions []Position `json:"positions"`
}