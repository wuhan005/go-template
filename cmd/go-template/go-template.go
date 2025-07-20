// Copyright 2025 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/wuhan005/go-template/internal/conf"
	"github.com/wuhan005/go-template/internal/db"
	"github.com/wuhan005/go-template/internal/route"
)

func main() {
	host := flag.String("host", "0.0.0.0", "host to listen")
	port := flag.Int("port", 8000, "port to listen")
	flag.Parse()

	if err := conf.Init(); err != nil {
		logrus.WithError(err).Fatal("Failed to initialize configuration")
	}

	db, err := db.Init()
	if err != nil {
		logrus.WithError(err).Fatal("Failed to initialize database")
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	c := make(chan os.Signal, 1)
	// Trigger graceful shutdown on SIGINT or SIGTERM.
	// The default signal sent by the `kill` command is SIGTERM,
	// which is taken as the graceful shutdown signal for many systems, e.g. Kubernetes, Gunicorn.
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	address := fmt.Sprintf("%s:%d", *host, *port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to listen")
	}

	f := route.New(db)
	server := http.Server{
		Handler:           f,
		ReadHeaderTimeout: 3 * time.Second,
	}

	go func() {
		<-c
		_ = server.Shutdown(ctx)
		cancel()
	}()

	logrus.WithField("address", address).Info("Server is running")
	if err := server.Serve(listener); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			logrus.WithError(err).Error("Failed to serve")
			cancel()
		}
	}

	// Wait for CTRL-C.
	<-ctx.Done()
}
