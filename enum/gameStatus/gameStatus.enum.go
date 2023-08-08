package gameStatusType

type GameStatusType string

const (
	WAITING    GameStatusType = "Waiting"
	START      GameStatusType = "Start"
	NORMAL_END GameStatusType = "NormalEnd"
)
