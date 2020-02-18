package testtool

import (
	"fmt"
	"net"
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

type logger interface {
	Logf(format string, args ...interface{})
}

// WaitForPort tries to connect to given local port with constant backoff.
//
// Returns error if port is not ready after timeout.
func WaitForPort(l logger, timeout time.Duration, port string) error {
	stopTimer := time.NewTimer(timeout)
	defer stopTimer.Stop()

	t := time.NewTicker(time.Millisecond * 100)
	defer t.Stop()

	for {
		select {
		case <-stopTimer.C:
			return fmt.Errorf("no server started listening on port %s after timeout %s", port, timeout)
		case <-t.C:
			if err := portIsReady(port); err != nil {
				l.Logf("waiting for port: %s\n", err)
				break
			}
			return nil
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
