package client

import (
	"os"
)

// 判断所给路径是否可用
func PathExist(_path string) bool {
	_, err := os.Stat(_path)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}
