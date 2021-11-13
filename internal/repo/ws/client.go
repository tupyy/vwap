package ws

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/tupyy/vwap/internal/entity"
	"github.com/tupyy/vwap/internal/log"
	"golang.org/x/net/websocket"
)

type WSClient struct {
	conn         *websocket.Conn
	tradingPairs []string
}

func NewClient(tradingPairs []string) *WSClient {
	return &WSClient{tradingPairs: tradingPairs}
}

func (c *WSClient) Connect(ctx context.Context, endpoint string) error {
	if c.conn != nil {
		return nil
	}

	doneCh := make(chan struct{}, 1)
	errCh := make(chan error, 1)

	go func() {
		conn, err := websocket.Dial(endpoint, "", "http://localhost/")
		if err != nil {
			errCh <- err

			return
		}

		c.conn = conn

		doneCh <- struct{}{}
	}()

	select {
	case <-time.After(10 * time.Second):
		return errors.New("timeout while connecting to wg")
	case <-ctx.Done():
		return ctx.Err()
	case errConn := <-errCh:
		return errConn
	case <-doneCh:
		log.GetLogger().Info("client connected to ws")
	}

	return nil
}

func (c *WSClient) Disconnect() error {
	return c.conn.Close()
}

func (c *WSClient) Receive(ctx context.Context, outputCh chan<- interface{}, errCh chan<- error) {
	logger := log.GetLogger()

	go func() {
		for {
			msg, err := readWs(c.conn)
			if err != nil {
				errCh <- err

				continue
			}

			logger.Trace("receive new message %s", msg.MessageType.String())
			switch msg.MessageType {
			case errorMessageType:
				logger.Error("received error message: %s", string(msg.Message))
			case tickerMessageType:
				var t entity.Ticker
				err := json.Unmarshal(msg.Message, &t)
				if err != nil {
					logger.Error("%+v", err)
				}

				outputCh <- t
			case heartBeatMessageType:
				var t entity.HeartBeat
				err := json.Unmarshal(msg.Message, &t)
				if err != nil {
					logger.Error("%+v", err)
				}

				outputCh <- t
			}
		}
	}()
}

func (c *WSClient) Subscribe(ctx context.Context) error {
	msg := subscribeMessage{
		MessageType: "subscribe",
		ProductIDs:  c.tradingPairs,
		Channels:    []string{"heartbeat", "ticker"},
	}

	return c.makeSubcription(ctx, msg)

}

func (c *WSClient) Unsubscribe(ctx context.Context) error {
	msg := subscribeMessage{
		MessageType: "unsubscribe",
		ProductIDs:  c.tradingPairs,
		Channels:    []string{"heartbeat", "ticker"},
	}

	return c.makeSubcription(ctx, msg)

}

func (c *WSClient) makeSubcription(ctx context.Context, msg subscribeMessage) error {
	b, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	err = writeToWs(c.conn, b)
	if err != nil {
		return err
	}

	msgCh := make(chan message, 1)
	errCh := make(chan error, 1)
	go func() {
		msg, err := readWs(c.conn)
		if err != nil {
			errCh <- err

			return
		}

		msgCh <- msg
	}()

	select {
	case <-ctx.Done():
		return errors.New("context canceled")
	case <-time.After(10 * time.Second):
		return errors.New("timeout while reading subscribe answer")
	case err := <-errCh:
		return fmt.Errorf("error subscribing to %v: %v", c.tradingPairs, err)
	case m := <-msgCh:
		if m.MessageType == errorMessageType {
			return fmt.Errorf("error receiving subscribe answer message: %s", string(m.Message))
		}
	}

	return nil
}
