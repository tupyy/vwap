package ws

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"golang.org/x/net/websocket"

	"github.com/tupyy/vwap/internal/entity"
	"github.com/tupyy/vwap/internal/log"
)

type WSClient struct {
	// conn -- websocket connection
	conn *websocket.Conn
	// TradingPairs -- list of trading pairs
	tradingPairs []string
	// doneCh -- channel used to close the reader
	doneCh chan chan interface{}
}

func NewClient(tradingPairs []string) *WSClient {
	return &WSClient{
		tradingPairs: tradingPairs,
		doneCh:       make(chan chan interface{}, 1),
	}
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
		log.GetLogger().Infof("client connected to ws")
	}

	return nil
}

func (c *WSClient) Shutdown() error {
	if c.conn == nil {
		return nil
	}

	log.GetLogger().Debugf("closing receiver")

	// stop receiver
	retCh := make(chan interface{}, 1)
	c.doneCh <- retCh

	<-retCh
	log.GetLogger().Debugf("receiver closed")

	err := c.conn.Close()
	if err != nil {
		return err
	}

	c.conn = nil

	return err
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
