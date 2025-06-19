package server

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	std_path "path"
	"strconv"
	"trans/common"
)

func StartServer(addr string, path string) {
	ser := http.NewServeMux()
	ser.HandleFunc(common.ApiListFileUrl, handleListFiles(http.MethodGet, path))
	ser.HandleFunc(common.ApiDownloadUrl, handleDownloadFile(http.MethodGet, path))
	log.Printf("start server at %s", addr)
	if err := http.ListenAndServe(addr, ser); err != nil {
		common.ExitWithError(err)
	}
}

func handleListFiles(method string, path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		fileList, _ := common.GetFileList(path)
		data, err := json.Marshal(fileList)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.Write(data)
	}
}

func handleDownloadFile(method string, path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		filename := r.URL.Query().Get("filename")
		if filename == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		filepath := std_path.Join(path, filename)
		if !common.FileIsExist(filepath) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		file, err := os.Open(filepath)
		if err != nil {
			log.Printf("cannot open file %s: %v", filepath, err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		stat, err := file.Stat()
		if err != nil {
			log.Printf("cannot get file info: %v", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		defer file.Close()
		w.Header().Add("Content-Type", "application/octet-stream")
		w.Header().Add("Content-Disposition", "attachment; filename="+filename)
		w.Header().Add("Content-Length", strconv.Itoa(int(stat.Size())))
		if _, err = io.Copy(w, file); err != nil {
			log.Printf("read error: %v", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		log.Printf("send file %s success", filepath)
	}
}
