package tcp

import (
	"context"
	"fmt"
	"github.com/flandersRin/tcp/lib/logger"
	"github.com/flandersRin/tcp/pkg/interface"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// Config stores tcp server properties; like Address: ":6379"
type Config struct {
	Address string
}

func ListenAndServeWithSignal(cfg *Config, handler _interface.Handler) error {

	// To accept close signal by value is struct{}{}
	closeChan := make(chan struct{})
	// In order to receive Use Ctrl + C signal etc..
	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)

	// handle syscall signal; Graceful exit
	go func() {
		sig := <-sigCh
		switch sig {
		case syscall.SIGHUP, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT:
			closeChan <- struct{}{}
		}
	}()

	listener, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("bind: %s, start listening...", cfg.Address))
	ListenAndServe(listener, handler, closeChan)

	return nil
}

// ListenAndServe binds port and handle requests, blocking until close
func ListenAndServe(listener net.Listener, handler _interface.Handler, closeChan <-chan struct{}) {
	// listen signal
	go func() {
		<-closeChan
		logger.Info("serve shutting down...")
		_ = listener.Close() // listener.Accept() will return err immediately
		_ = handler.Close()  // close connection
	}()

	// Stop running normally
	defer func() {
		_ = listener.Close()
		_ = handler.Close()
	}()

	ctx := context.Background()
	var waitDone sync.WaitGroup
	for {
		conn, err := listener.Accept()
		if err != nil {
			break
		}
		// handle
		logger.Info("accept link")
		waitDone.Add(1)

		// goroutine concurrency accept and handle request
		go func() {
			defer func() {
				waitDone.Done()
			}()
			handler.Handle(ctx, conn)
		}()
	}
	waitDone.Wait()
}
