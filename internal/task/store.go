package task

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)


type Store struct {
	path string
}


func NewStore() (*Store, error) {
	dir, err := dataDir()
	if err != nil {
		return nil, err
	}


	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("creating data directory %s: %w", dir, err)
	}

	return &Store{path: filepath.Join(dir, "tasks.json")}, nil
}


func dataDir() (string, error) {

	if xdg := os.Getenv("XDG_DATA_HOME"); xdg != "" {
		return filepath.Join(xdg, "tascii"), nil
	}


	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not determine home directory: %w", err)
	}

	return filepath.Join(home, ".local", "share", "tascii"), nil
}


func (s *Store) Load() ([]Task, error) {
	data, err := os.ReadFile(s.path)


	if os.IsNotExist(err) {
		return []Task{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", s.path, err)
	}


	var tasks []Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		return nil, fmt.Errorf("parsing tasks file (may be corrupted): %w", err)
	}

	return tasks, nil
}


func (s *Store) Save(tasks []Task) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return fmt.Errorf("encoding tasks: %w", err)
	}

	//read/write for owner, readonly for others
	return os.WriteFile(s.path, data, 0644)
}





func (s *Store) NextID(tasks []Task) int {
	max := 0
	for _, t := range tasks {
		if t.ID > max {
			max = t.ID
		}
	}
	return max + 1
}


func (s *Store) Path() string {
	return s.path
}
