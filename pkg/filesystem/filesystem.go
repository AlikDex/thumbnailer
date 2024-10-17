package filesystem

import (
	"os"
)

func IsFile(pathname string) bool {
	info, err := os.Stat(pathname)

	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}

func IsDir(dirname string) bool {
	info, err := os.Stat(dirname)

	if os.IsNotExist(err) {
		return false
	}

	return info.IsDir()
}

func IsExists(targetPath string) bool {
	_, err := os.Stat(targetPath)

	return !os.IsNotExist(err)
}
