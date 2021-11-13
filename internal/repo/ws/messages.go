// this messages are internal to coinbase ws client and they are not exposed.
package ws

// messageType defines the type of the message received from coinbase
type messageType int

const (
	subscribeMessageType messageType = iota
	subcribingMessageType
	unsubscribeMessageType
	errorMessageType
	tickerMessageType
	heartBeatMessageType
	unknownMessageType
)

func (m messageType) String() string {
	switch m {
	case subscribeMessageType:
		return "subscribe message"
	case subcribingMessageType:
		return "subscribing message"
	case unknownMessageType:
		return "unsubscribe message"
	case errorMessageType:
		return "error message"
	case tickerMessageType:
		return "ticker message"
	case heartBeatMessageType:
		return "heartbeat message"
	default:
		return "unknown message"
	}
}

// message is a general type of message. It holds the type of message and the actual message as byte.
type message struct {
	MessageType messageType
	Message     []byte
}

type errorMessage struct {
	MessageType string `json:"type"`
	Message     string `json:"message"`
	Reason      string `json:"reason"`
}

// message -- subscribe / unsubscribe message
type subscribeMessage struct {
	MessageType string   `json:"type"`
	ProductIDs  []string `json:"product_ids"`
	Channels    []string `json:"channels"`
}

type channels struct {
	Name       string   `json:"name"`
	ProductIds []string `json:"product_ids"`
}

type subscribeAnswerMessage struct {
	MessageType string     `json:"type"`
	Channels    []channels `json:"channels"`
}
