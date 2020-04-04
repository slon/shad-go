package main

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gitlab.com/slon/shad-go/tools/testtool"
)

const importPath = "gitlab.com/slon/shad-go/wscat"

var binCache testtool.BinCache

func TestMain(m *testing.M) {
	os.Exit(func() int {
		var teardown testtool.CloseFunc
		binCache, teardown = testtool.NewBinCache()
		defer teardown()

		return m.Run()
	}())
}

type Conn struct {
	in  io.WriteCloser
	out *bytes.Buffer
}

func startCommand(t *testing.T, addr string) (conn *Conn, stop func()) {
	t.Helper()

	binary, err := binCache.GetBinary(importPath)
	require.NoError(t, err)

	cmd := exec.Command(binary, "-addr", addr)
	cmd.Stderr = os.Stderr

	out := &bytes.Buffer{}
	cmd.Stdout = out

	stdin, err := cmd.StdinPipe()
	require.NoError(t, err)

	require.NoError(t, cmd.Start())

	conn = &Conn{
		in:  stdin,
		out: out,
	}

	done := make(chan struct{})
	go func() {
		assert.NoError(t, cmd.Wait())
		close(done)
	}()

	stop = func() {
		defer func() {
			_ = cmd.Process.Kill()
			<-done
		}()

		// try killing softly
		_ = cmd.Process.Signal(syscall.SIGTERM)

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
		defer cancel()

		select {
		case <-done:
		case <-ctx.Done():
			t.Fatalf("client shutdown timed out")
		}
	}

	return
}

func TestWScat(t *testing.T) {
	upgrader := websocket.Upgrader{}

	var received, sent [][]byte
	h := func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		require.NoError(t, err)
		defer func() { _ = c.Close() }()

		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				t.Logf("error reading message: %s", err)
				break
			}
			received = append(received, message)

			resp := []byte(testtool.RandomName())
			err = c.WriteMessage(websocket.TextMessage, resp)
			require.NoError(t, err)
			sent = append(sent, resp)
		}
	}

	s := httptest.NewServer(http.HandlerFunc(h))
	defer s.Close()

	wsURL := strings.Replace(s.URL, "http", "ws", 1)
	t.Logf("starting ws server %s", wsURL)

	conn, stop := startCommand(t, wsURL)
	defer stop()

	var in [][]byte
	for i := 0; i < 100; i++ {
		msg := []byte(testtool.RandomName())
		in = append(in, msg)

		_, err := conn.in.Write(append(msg, '\n'))
		require.NoError(t, err)
	}

	// give client time to make a request
	time.Sleep(time.Millisecond * 100)
	stop()

	require.Equal(t, bytes.Join(in, nil), bytes.Join(received, nil))
	require.Equal(t, bytes.Join(sent, nil), conn.out.Bytes())
}
