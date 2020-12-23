package utils

import (
	"os"
	"path/filepath"
)

func IsExist(elm ...string) bool {
	path := filepath.Join(elm...)
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		return true
	}

	f.Close()
	return true
}
