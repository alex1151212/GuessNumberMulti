package playerStatusType

type PlayerStatusType string

const (
	INLOBBY       PlayerStatusType = "InLobby"
	WAITING_START PlayerStatusType = "WaitingStart"
	PLAYING       PlayerStatusType = "Playing"
)
