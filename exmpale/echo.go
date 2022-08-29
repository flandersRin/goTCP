package exmpale

import (
	"bufio"
	"context"
	"github.com/flandersRin/tcp/lib/logger"
	"github.com/flandersRin/tcp/lib/sync/atomic"
	"github.com/flandersRin/tcp/lib/sync/wait"
	"io"
	"net"
	"sync"
	"time"
)

type EchoHandler struct {
	activeConn sync.Map       // Add many connection into safe sync.Map
	closing    atomic.Boolean // concurrency safe set closing mark
}

func NewHandler() *EchoHandler {
	return &EchoHandler{}
}

// EchoClient is client for EchoHandler, using for testing
type EchoClient struct {
	Conn    net.Conn
	Waiting wait.Wait
}

// Close closed connection
func (c *EchoClient) Close() error {
	c.Waiting.WaitWithTimeout(10 * time.Second)
	c.Conn.Close()
	return nil
}

// Handle echos received line to client
func (h *EchoHandler) Handle(ctx context.Context, conn net.Conn) {
	if h.closing.Get() {
		// closing handler refuse new connection
		_ = conn.Close()
	}

	client := &EchoClient{
		Conn: conn,
	}
	// make the hashMap to hashSet
	h.activeConn.Store(client, struct{}{})

	reader := bufio.NewReader(conn)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				logger.Info("connection close")
				h.activeConn.Delete(client)
			} else {
				logger.Warn(err)
			}
			return
		}

		client.Waiting.Add(1)
		_, _ = conn.Write([]byte(msg))
		client.Waiting.Done()
	}

}

func (h *EchoHandler) Close() error {
	logger.Info("handler shutting down...")
	h.closing.Set(true)
	h.activeConn.Range(func(key, value any) bool {
		client := key.(*EchoClient)
		_ = client.Close()
		return true // Chain next
	})
	return nil
}
