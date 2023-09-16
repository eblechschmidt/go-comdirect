package state

import (
	"fmt"
	"os"

	"github.com/adrg/xdg"
	"gopkg.in/yaml.v3"
)

type state struct {
	Since map[string]string `yaml:"since"`
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func new() *state {
	s := state{}
	s.Since = make(map[string]string)
	return &s
}

func readState() (*state, error) {
	s := new()
	stateFile, err := xdg.StateFile("comdirect/state.yaml")
	if err != nil {
		return nil, fmt.Errorf("could not get state file: %w", err)
	}
	if !fileExists(stateFile) {
		return s, nil
	}
	yamlFile, err := os.ReadFile(stateFile)
	if err != nil {
		return nil, fmt.Errorf("could not read state file: %w", err)
	}
	err = yaml.Unmarshal(yamlFile, &s)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal state file: %w", err)
	}
	return s, nil
}

func saveState(s *state) error {
	stateFile, err := xdg.StateFile("comdirect/state.yaml")
	if err != nil {
		return fmt.Errorf("could not get state file: %w", err)
	}
	d, err := yaml.Marshal(s)
	if err != nil {
		return err
	}

	return os.WriteFile(stateFile, d, 0666)
}

func GetLastSince(accountID string) (string, error) {
	s, err := readState()
	if err != nil {
		return "", err
	}
	if since, ok := s.Since[accountID]; ok {
		return since, nil
	}
	return "", nil
}

func SetLastSince(accountID, since string) error {
	s, err := readState()
	if err != nil {
		return err
	}
	s.Since[accountID] = since

	saveState(s)
	return nil
}
