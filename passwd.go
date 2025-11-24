package main

import (
	"strings"
)

type PasswdEntry struct {
	Username string
	Password string
	UID      string
	GID      string
	GECOS    string
	HomeDir  string
	Shell    string
}

func ParsePasswdEntry(entry string) *PasswdEntry {
	fields := strings.Split(entry, ":")
	return &PasswdEntry{
		Username: fields[0],
		Password: fields[1],
		UID:      fields[2],
		GID:      fields[3],
		GECOS:    fields[4],
		HomeDir:  fields[5],
		Shell:    fields[6],
	}
}
