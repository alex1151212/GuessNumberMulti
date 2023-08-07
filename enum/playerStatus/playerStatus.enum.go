package playerStatusType

type PlayerStatusType string

const (
	INLOBBY      PlayerStatusType = "InLobby"
	WAITINGSTART PlayerStatusType = "WaitingStart"
	PLAYING      PlayerStatusType = "Playing"
)
