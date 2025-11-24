package main

import (
	"os"
	"strings"
)

type LinuxUser struct {
	Username string
	Password string
	UID      string
	GID      string
	GECOS    string
	HomeDir  string
	Shell    string
}

func ParseLinuxUser(entry string) *LinuxUser {
	fields := strings.Split(entry, ":")
	return &LinuxUser{
		Username: fields[0],
		Password: fields[1],
		UID:      fields[2],
		GID:      fields[3],
		GECOS:    fields[4],
		HomeDir:  fields[5],
		Shell:    fields[6],
	}
}

func listManagedLinuxUsers() (map[string]*LinuxUser, error) {
	passwd, err := os.ReadFile("/etc/passwd")
	if err != nil {
		return nil, err
	}
	users := make(map[string]*LinuxUser, 0)
	for entry := range strings.SplitSeq(string(passwd), "\n") {
		fields := strings.Split(entry, ":")
		if len(fields) != 7 {
			continue
		}
		linuxUser := ParseLinuxUser(entry)
		switch linuxUser.GECOS {
		case _gecosInfoGoDCVManaged:
			users[linuxUser.Username] = linuxUser
		}
	}
	return users, nil
}
