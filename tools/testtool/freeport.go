package testtool

import (
	"context"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

// GetFreePort returns free local tcp port.
func GetFreePort() (string, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return "", err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return "", err
	}
	defer func() { _ = l.Close() }()

	p := l.Addr().(*net.TCPAddr).Port

	return strconv.Itoa(p), nil
}

// WaitForPort tries to connect to given local port with constant backoff.
//
// Can be canceled via ctx.
func WaitForPort(ctx context.Context, port string) {
	t := time.NewTicker(time.Millisecond * 100)
	defer t.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			if err := portIsReady(port); err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "waiting for port: %s\n", err)
				break
			}
			return
		}
	}
}

func portIsReady(port string) error {
	conn, err := net.Dial("tcp", net.JoinHostPort("localhost", port))
	if err != nil {
		return err
	}
	return conn.Close()
}
