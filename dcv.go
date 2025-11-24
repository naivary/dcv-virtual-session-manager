package main

import (
	"encoding/json"
	"os/exec"
)

type DCVSession struct {
	Name        string
	User        string
	Owner       string
	StorageRoot string
	ID          string
}

func createVirtualSessionFromPasswdEntry(e *PasswdEntry, storageRoot string) error {
	session := DCVSession{
		Name:        e.Username,
		User:        e.Username,
		Owner:       e.Username,
		ID:          e.Username,
		StorageRoot: storageRoot,
	}
	return createVirtualSession(&session)
}

func createVirtualSession(s *DCVSession) error {
	isAlreadyCreated, err := isVirtualSessionAlreadyCreated(s.ID)
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

func isVirtualSessionAlreadyCreated(id string) (bool, error) {
	var ids []struct {
		ID string `json:"id"`
	}
	cmd := exec.Command(
		"dcv",
		"list-sessions",
		"--type", "virtual",
		"--json",
	)
	data, err := cmd.Output()
	if err != nil {
		return true, err
	}
	err = json.Unmarshal(data, &ids)
	if err != nil {
		return true, err
	}
	for _, i := range ids {
		if i.ID == id {
			return true, nil
		}
	}
	return false, nil
}
