package fsx

import "os"

func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func IsDirectory(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.Mode().IsDir()
}

func IsFile(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.Mode().IsRegular()
}
