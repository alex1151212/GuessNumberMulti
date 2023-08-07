package messageType

type MessageType string

const (
	INIT         MessageType = "Init"
	PLAYING      MessageType = "Playing"
	GET_GAMES    MessageType = "GetGames"
	CREATE_GAMES MessageType = "CreateGames"
	GET_PLAYERS  MessageType = "GetPlayers"
	JOIN_GAME    MessageType = "JoinGame"
	GAME_START   MessageType = "GameStart"
)
