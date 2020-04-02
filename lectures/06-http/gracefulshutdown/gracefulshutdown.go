package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("pong"))
}

func run() error {
	srv := &http.Server{
		Addr:    ":8080",
		Handler: http.HandlerFunc(handler),
	}
	serveChan := make(chan error, 1)
	go func() {
		serveChan <- srv.ListenAndServe()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-stop:
		fmt.Println("shutting down gracefully")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return srv.Shutdown(ctx)

	case err := <-serveChan:
		return err
	}
}
