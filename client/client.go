package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	std_path "path"
	"strconv"
	"time"
	"trans/common"
)

func ListFiles(addr string) {
	files := listFiles(addr)
	filenames := make([]string, len(files))
	for i, file := range files {
		filenames[i] = file.Name
	}
	fmt.Println("files on remote server:")
	for i, file := range files {
		size := common.GetHumanizedSize(file.Size)
		fmt.Printf("%3d  %6s  %s\n", i+1, size, file.Name)
	}
}

func listFiles(addr string) common.FileList {
	url := fmt.Sprintf("http://%s%s", addr, common.ApiListFileUrl)
	resp, err := http.Get(url)
	if err != nil {
		common.ExitWithError(errors.Unwrap(err))
	}
	defer resp.Body.Close()
	var res common.FileList
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		common.ExitWithError(err)
	}
	return res
}

func DownloadFile(addr string, path string, filenames []string) {
	for i, filename := range filenames {
		downloadFile(addr, path, filename)
		if i < len(filenames)-1 {
			fmt.Println("---")
		}
	}
}

func DownloadAllFile(addr string, path string) {
	files := listFiles(addr)
	filenames := make([]string, len(files))
	for i, file := range files {
		filenames[i] = file.Name
	}
	DownloadFile(addr, path, filenames)
}

func downloadFile(addr string, path string, filename string) {
	log.Printf("downloading %s", filename)
	url := fmt.Sprintf("http://%s%s?filename=%s", addr, common.ApiDownloadUrl, filename)
	resp, err := http.Get(url)
	if err != nil {
		log.Println(errors.Unwrap(err).Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("download file failed: %v", filename)
		return
	}

	filepath := std_path.Join(path, filename)
	file, err := os.Create(filepath)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer file.Close()

	contentLenght, _ := strconv.Atoi(resp.Header.Get("Content-Length"))
	buf := make([]byte, 4096)

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	var received int
	go printProgress(&received, &contentLenght, ticker)

	for received = 0; received < contentLenght; {
		n, err := resp.Body.Read(buf)
		if err != nil && err != io.EOF {
			log.Printf("download file failed: %v %v", filename, err)
			os.Remove(filepath)
			break
		}
		file.Write(buf[:n])
		received += n
	}
	printProgress(&received, &contentLenght, nil)
	log.Printf("downloaded file to %s", filepath)
}

func printProgress(received *int, total *int, ticker *time.Ticker) {
	if ticker != nil {
		for range ticker.C {
			fmt.Printf(
				"\r%s/%s",
				common.GetHumanizedSize(int64(*received)),
				common.GetHumanizedSize(int64(*total)),
			)
		}
	}
	fmt.Printf(
		"\r%s/%s\n",
		common.GetHumanizedSize(int64(*received)),
		common.GetHumanizedSize(int64(*total)),
	)
}
