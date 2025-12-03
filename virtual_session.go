package main

import (
	"encoding/json"
	"os/exec"
)

type VirtualSession struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Owner       string `json:"owner"`
	User        string `json:"user"`
	StorageRoot string `json:"storage-root"`
}

func pruneVirtualSessions() error {
	users, err := listManagedLinuxUsers()
	if err != nil {
		return err
	}
	sessions, err := listVirtualSessions()
	if err != nil {
		return err
	}
	for _, session := range sessions {
		_, ok := users[session.Owner]
		// user and session are existing
		if ok {
			continue
		}
		// session is existing without the user
		err := deleteVirtualSession(session.ID)
		if err != nil {
			return err
		}
	}
	return nil
}

func createVirtualSessionFromLinuxUser(e *LinuxUser, storageRoot string) error {
	session := VirtualSession{
		Name:        e.Username,
		User:        e.Username,
		Owner:       e.Username,
		ID:          e.Username,
		StorageRoot: storageRoot,
	}
	return createVirtualSession(&session)
}

func createVirtualSessionFromPasswd(storagePath string) error {
	users, err := listManagedLinuxUsers()
	if err != nil {
		return err
	}
	for _, user := range users {
		err := createVirtualSessionFromLinuxUser(user, storagePath)
		if err != nil {
			return err
		}
	}
	return nil
}

func createVirtualSession(s *VirtualSession) error {
	isAlreadyCreated, err := isVirtualSessionCreated(s.ID)
	if err != nil {
		return err
	}
	if isAlreadyCreated {
		return nil
	}
	cmd := exec.Command(
		"dcv",
		"create-session",
		"--type", "virtual",
		"--name", s.Name,
		"--user", s.User,
		"--owner", s.Owner,
		"--storage-root", s.StorageRoot,
		s.ID,
	)
	return cmd.Run()
}

func deleteVirtualSession(id string) error {
	cmd := exec.Command(
		"dcv",
		"close-session",
		id,
	)
	return cmd.Run()
}

// returns the ids of the current virtual sessions
func listVirtualSessions() ([]VirtualSession, error) {
	var sessions []VirtualSession
	cmd := exec.Command(
		"dcv",
		"list-sessions",
		"--type", "virtual",
		"--json",
	)
	data, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &sessions)
	if err != nil {
		return nil, err
	}
	return sessions, nil
}

func isVirtualSessionCreated(id string) (bool, error) {
	sessions, err := listVirtualSessions()
	if err != nil {
		return true, err
	}
	for _, s := range sessions {
		if s.ID == id {
			return true, nil
		}
	}
	return false, nil
}
