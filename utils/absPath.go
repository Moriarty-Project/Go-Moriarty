package utils

import (
	"fmt"
	"os"
	"path"
	"strings"
)

// gets the absolute filepath to where ever you need. Attempts multiple things, but if none come back, returns the final error
// mostly is trying to find the correct file, returning an absolute path is purely by chance
func GetAbsolutePath(to string) (string, error) {
	_, err := os.Stat(to)
	if err == nil {
		return to, nil
	}
	// that didn't work
	// next, try adding the cwd.
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	_, err = os.Stat(path.Join(cwd, to))
	if err == nil {
		return path.Join(cwd, to), nil
	}
	// I guess we'll try merging them with ..'s until it fits?
	middle := "../"
	for i := 0; i < 5 && i < strings.Count(to, "/"); i++ {
		_, err = os.Stat(path.Join(cwd, middle, to))
		if err == nil {
			return path.Join(cwd, middle, to), nil
		}
		middle = middle + middle
	}
	return "", fmt.Errorf("could not find any path to %v", to)
}
