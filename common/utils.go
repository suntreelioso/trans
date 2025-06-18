package common

import (
	"fmt"
	"os"
	"strings"
)

func ExitWithError(err error) {
	fmt.Printf("error: %v\n", err.Error())
	os.Exit(1)
}

func DirIsExist(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func FileIsExist(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.Mode().IsRegular()
}

func GetCwd() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}
	return dir
}

func GetFileList(path string) (FileList, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	fileList := make(FileList, 0)
	for _, entry := range entries {
		if !entry.Type().IsRegular() {
			continue
		}
		if IsHiddenFile(entry.Name()) {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			return nil, err
		}
		fileInfo := FileInfo{
			Name: entry.Name(),
			Size: info.Size(),
		}
		fileList = append(fileList, fileInfo)
	}
	return fileList, nil
}

func GetHumanizedSize(size int64) string {
	const n = 1000
	if size < n {
		return fmt.Sprintf("%dB", size)
	} else if size < n*n {
		return fmt.Sprintf("%.1fK", float64(size)/n)
	} else if size < n*n*n {
		return fmt.Sprintf("%.1fM", float64(size)/n/n)
	} else {
		return fmt.Sprintf("%.1fG", float64(size)/n/n/n)
	}
}

func IsHiddenFile(fileName string) bool {
	return strings.HasPrefix(fileName, ".")
}
