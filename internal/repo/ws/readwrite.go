package ws

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/tupyy/vwap/internal/log"
)

func readWs(r io.Reader, buffer []byte) (message, error) {
	logger := log.GetLogger()

	n, err := r.Read(buffer)
	if err != nil {
		return message{}, err
	}

	logger.Tracef("read %d bytes. read msg from websocket %s", n, string(buffer[:n]))

	msgData := make([]byte, n)
	copy(msgData, buffer[:n])

	msgType, err := getMessageType(msgData)
	if err != nil {
		return message{}, err
	}

	return message{
		MessageType: msgType,
		Message:     msgData,
	}, nil
}

func writeToWs(w io.Writer, msg []byte) error {
	log.GetLogger().Tracef("write message to ws: %s", string(msg))

	n, err := w.Write(msg)
	if err != nil {
		return err
	}

	log.GetLogger().Debugf("wrote %d bytes to websocket", n)

	return nil
}

func getMessageType(msg []byte) (messageType, error) {
	m := struct {
		MsgType string `json:"type"`
	}{}

	err := json.Unmarshal(msg, &m)
	if err != nil {
		return 0, fmt.Errorf("error reading msg type: %+v. message: %s", err, string(msg))
	}

	switch m.MsgType {
	case "subscriptions":
		return subcribingMessageType, nil
	case "error":
		return errorMessageType, nil
	case "ticker":
		return tickerMessageType, nil
	case "heartbeat":
		return heartBeatMessageType, nil
	default:
		return unknownMessageType, nil
	}
}
