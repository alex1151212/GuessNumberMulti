package gameStatusType

type GameStatusType string

const (
	WAITING    GameStatusType = "Waiting"
	START      GameStatusType = "Start"
	NORMAL_END GameStatusType = "NormalEnd"
	EARLY_END  GameStatusType = "EarlyEnd"
)
