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
	for _, filename := range filenames {
		log.Printf("downloading %s", filename)
		url := fmt.Sprintf("http://%s%s?filename=%s", addr, common.ApiDownloadUrl, filename)
		resp, err := http.Get(url)
		if err != nil {
			log.Println(errors.Unwrap(err).Error())
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Printf("download file failed: %v", filename)
			continue
		}

		filepath := std_path.Join(path, filename)
		file, err := os.Create(filepath)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		if _, err = io.Copy(file, resp.Body); err != nil {
			log.Printf("download file failed: %v %v", filename, err.Error())
			file.Close()
			os.Remove(filepath)
			continue
		}
		file.Close()
		resp.Body.Close()
		log.Printf("downloaded file to %s", filepath)
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
