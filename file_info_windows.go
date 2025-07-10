//go:build windows

package main

import (
	"os"
	"syscall"
	"time"
)

func getFileCreationTime(path string) (time.Time, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return time.Time{}, err
	}
	if stat, ok := fileInfo.Sys().(*syscall.Win32FileAttributeData); ok {
		return time.Unix(0, stat.CreationTime.Nanoseconds()), nil
	}
	return fileInfo.ModTime(), nil
}
