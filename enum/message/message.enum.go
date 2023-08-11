package messageType

type MessageType string

const (
	INIT             MessageType = "Init"
	PLAYING          MessageType = "Playing"
	GET_GAMES        MessageType = "GetGames"
	CREATE_GAME      MessageType = "CreateGame"
	GET_PLAYERS      MessageType = "GetPlayers"
	JOIN_GAME        MessageType = "JoinGame"
	DELETE_GAME      MessageType = "DeleteGame"
	LEAVE_GAME       MessageType = "LeaveGame"
	INPUT_GAMEANSWER MessageType = "InputGameAnswer"

	GAME_START MessageType = "GameStart"
	NORMAL_END MessageType = "NormalEnd"

	ERROR MessageType = "Error"
)
