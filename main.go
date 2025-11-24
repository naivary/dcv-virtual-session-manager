package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"time"
)

const _gecosInfoGoDCVManaged = "go_dcv_managed"

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(context.Context) error {
	var storagePath string
	var period time.Duration
	flag.StringVar(&storagePath, "storagepath", "/opt/session-storage", "path to use for --storage-root")
	flag.DurationVar(&period, "period", 5*time.Minute, "period at which /etc/passwd is re-read")
	flag.Parse()

	err := createVirtualSessionFromPasswd(storagePath)
	if err != nil {
		return err
	}
	ticker := time.NewTicker(period)
	defer ticker.Stop()
	slog.Info("Re-reading /etc/passwd periodically", "period", period)
	for range ticker.C {
		// remove all sessions for which the user does not exist anymore
		err = pruneVirtualSessions()
		if err != nil {
			return err
		}
		// create new virtual sessions for newly added sessions
		err = createVirtualSessionFromPasswd(storagePath)
		if err != nil {
			slog.Error("failed to create virtual sessions from passwd", "msg", err.Error())
		}
	}
	<-make(chan struct{})
	return nil
}
