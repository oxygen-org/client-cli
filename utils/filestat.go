package utils

import "os"

// FileExists 判断文件是否存在
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// DirExists 判断dir是否存在
func DirExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// CreateDirIfNotExists x
func CreateDirIfNotExists(dirpath string){
	os.MkdirAll(dirpath, os.ModePerm)
}

