package messageType

type MessageType string

const (
	INIT         MessageType = "Init"
	PLAYING      MessageType = "Playing"
	GET_GAMES    MessageType = "GetGames"
	CREATE_GAMES MessageType = "CreateGames"
	GET_PLAYERS  MessageType = "GetPlayers"
	JOIN_GAME    MessageType = "JoinGame"
	DELETE_GAME  MessageType = "DeleteGame"

	GAME_START MessageType = "GameStart"
	NORMAL_END MessageType = "NormalEnd"

	ERROR MessageType = "Error"
)
