//go:build !windows

package main

import (
	"os"
	"time"
)

func getFileCreationTime(path string) (time.Time, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return time.Time{}, err
	}
	// On non-Windows systems, we fall back to modification time.
	return fileInfo.ModTime(), nil
}
