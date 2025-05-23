package backup

import (
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/v3lmx/counter/internal/core"
)

type FileBackup struct {
	Current *os.File
	Best    *os.File
}

func NewFileBackup(currentPath string, bestPath string) (FileBackup, error) {
	fb := FileBackup{}

	current, err := os.OpenFile(currentPath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return fb, fmt.Errorf("Could not open/create current backup file: %v", err)
	}
	fb.Current = current

	best, err := os.OpenFile(bestPath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return fb, fmt.Errorf("Could not open/create best backup file: %v", err)
	}
	fb.Best = best

	return fb, nil
}

func (fb FileBackup) Backup(current uint64, best core.Best) error {
	err := fb.Current.Truncate(0)
	if err != nil {
		return fmt.Errorf("Could not truncate current backup: %v", err)
	}
	_, err = fb.Current.Seek(0, 0)
	if err != nil {
		return fmt.Errorf("Could not seek current backup: %v", err)
	}
	_, err = fb.Current.WriteString(strconv.Itoa(int(current)))
	if err != nil {
		return fmt.Errorf("Could not write current to backup: %v", err)
	}

	err = fb.Best.Truncate(0)
	if err != nil {
		return fmt.Errorf("Could not truncate best backup: %v", err)
	}
	_, err = fb.Best.Seek(0, 0)
	if err != nil {
		return fmt.Errorf("Could not seek best backup: %v", err)
	}
	_, err = fb.Best.WriteString(best.Format())
	if err != nil {
		return fmt.Errorf("Could not write best to backup: %v", err)
	}

	return nil
}

func (fb FileBackup) Recover() (uint64, core.Best, error) {
	currentBuffer, err := io.ReadAll(fb.Current)
	if err != nil {
		return 0, core.Best{}, fmt.Errorf("Could not recover current backup: %v", err)
	}

	currentStr := string(currentBuffer)
	if currentStr == "" {
		return 0, core.Best{}, nil
	}

	current, err := strconv.Atoi(string(currentBuffer))
	if err != nil {
		return 0, core.Best{}, fmt.Errorf("Could not parse current backup: %v", err)
	}

	bestBuffer, err := io.ReadAll(fb.Best)
	if err != nil {
		return 0, core.Best{}, fmt.Errorf("Could not recover best backup: %v", err)
	}

	bestStr := string(bestBuffer)
	if bestStr == "" {
		return 0, core.Best{}, nil
	}

	best, err := (core.ParseBest(bestStr))
	if err != nil {
		return 0, core.Best{}, fmt.Errorf("Could not parse current backup: %v", err)
	}

	return uint64(current), best, nil
}
