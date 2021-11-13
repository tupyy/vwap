package ws

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/tupyy/vwap/internal/entity"
	"github.com/tupyy/vwap/internal/log"
)

type WSClient struct {
	// conn -- websocket connection
	conn io.ReadWriter
	// TradingPairs -- list of trading pairs
	tradingPairs []string
	// doneCh -- channel used to close the reader
	doneCh chan chan interface{}
}

func NewClient(conn io.ReadWriter, tradingPairs []string) *WSClient {
	return &WSClient{
		conn:         conn,
		tradingPairs: tradingPairs,
		doneCh:       make(chan chan interface{}, 1),
	}
}

func (c *WSClient) Shutdown() {
	log.GetLogger().Debugf("closing receiver")

	// stop receiver
	retCh := make(chan interface{}, 1)
	c.doneCh <- retCh

	<-retCh
	log.GetLogger().Debugf("receiver closed")
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

			logger.Tracef("receive new message %s", msg.MessageType.String())
			switch msg.MessageType {
			case errorMessageType:
				errCh <- fmt.Errorf("received an error message: %s", string(msg.Message))
			case tickerMessageType:
				var t entity.Ticker
				err := json.Unmarshal(msg.Message, &t)
				if err != nil {
					logger.Errorf("%+v", err)
				}

				outputCh <- t
			case heartBeatMessageType:
				var t entity.HeartBeat
				err := json.Unmarshal(msg.Message, &t)
				if err != nil {
					logger.Errorf("%+v", err)
				}

				outputCh <- t

			}

			select {
			case <-ctx.Done():
				logger.Errorf("context canceled: %+v", ctx.Err())
				return
			case retCh := <-c.doneCh:
				retCh <- struct{}{}
				return
			default:
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
