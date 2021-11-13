package ws

import (
	"context"
	"errors"
	"time"

	"golang.org/x/net/websocket"
)

func Connect(ctx context.Context, endpoint string) (*websocket.Conn, error) {
	doneCh := make(chan *websocket.Conn, 1)
	errCh := make(chan error, 1)

	go func() {
		conn, err := websocket.Dial(endpoint, "", "http://localhost/")
		if err != nil {
			errCh <- err

			return
		}

		doneCh <- conn
	}()

	select {
	case <-time.After(10 * time.Second):
		return nil, errors.New("timeout while connecting to wg")
	case <-ctx.Done():
		return nil, ctx.Err()
	case errConn := <-errCh:
		return nil, errConn
	case c := <-doneCh:
		return c, nil
	}
}
