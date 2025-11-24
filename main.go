package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"
)

const _gecosInfoCreateVirtualSession = "dcv_create_virtual_session"

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
		err = createVirtualSessionFromPasswd(storagePath)
		if err != nil {
			slog.Error("failed to create virtual sessions from passwd", "msg", err.Error())
		}
	}
	<-make(chan struct{})
	return nil
}

func createVirtualSessionFromPasswd(storagePath string) error {
	passwd, err := os.ReadFile("/etc/passwd")
	if err != nil {
		return err
	}
	for entry := range strings.SplitSeq(string(passwd), "\n") {
		fields := strings.Split(entry, ":")
		if len(fields) != 7 {
			continue
		}
		passwdEntry := ParsePasswdEntry(entry)
		switch passwdEntry.GECOS {
		case _gecosInfoCreateVirtualSession:
			slog.Info("creating virtual session", "user", passwdEntry.Username)
			err := createVirtualSessionFromPasswdEntry(passwdEntry, storagePath)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
