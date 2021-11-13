package ws

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tupyy/vwap/internal/entity"
)

func TestRead(t *testing.T) {
	heartbeat := `{
    "type": "heartbeat",
    "sequence": 90,
    "last_trade_id": 20,
    "product_id": "BTC-USD",
    "time": "2014-11-07T08:19:28.464459Z"
}`

	r := bytes.NewBufferString(heartbeat)

	receivedMsg, err := readWs(r)
	assert.Nil(t, err, "err should be nil")

	assert.Equal(t, heartBeatMessageType, receivedMsg.MessageType)

	var ticker entity.Ticker
	err = json.Unmarshal(receivedMsg.Message, &ticker)
	assert.Nil(t, err)

	assert.Equal(t, int64(90), ticker.Sequence, "seq should be 90")
}

func TestReadUnkownMessage(t *testing.T) {
	unknown := `{
    "type": "unknown",
    "sequence": 90,
    "last_trade_id": 20,
    "product_id": "BTC-USD",
    "time": "2014-11-07T08:19:28.464459Z"
}`

	r := bytes.NewBufferString(unknown)

	receivedMsg, err := readWs(r)
	assert.Nil(t, err, "err should be nil")

	assert.Equal(t, unknownMessageType, receivedMsg.MessageType)
}

func TestWriter(t *testing.T) {
	msg := bytes.NewBufferString("hey").Bytes()

	var w bytes.Buffer

	err := writeToWs(&w, msg)
	assert.Nil(t, err, "should be nil")

	assert.Equal(t, "hey", w.String())
}
