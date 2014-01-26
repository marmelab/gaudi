package util

import (
	"os"
)

func IsDir(path string) bool {
	stat := getFileStat(path)
	if stat == nil {
		return false
	}

	return stat.Mode().IsDir()
}

func IsFile(path string) bool {
	stat := getFileStat(path)
	if stat == nil {
		return false
	}

	return stat.Mode().IsRegular()
}

func getFileStat(path string) os.FileInfo {
	f, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		return nil
	}

	return stat
}
