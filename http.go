package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	PathList      = "/files"
	PathDownload  = "/download/"
	QueryParamKey = "secret"
)

func runHttpServer(port int) {
	fn := "runHttpServer"
	http.HandleFunc(PathList, secretFilterHandler(fileListHandler))
	http.HandleFunc(PathDownload, secretFilterHandler(fileDownloadHandler))

	fmt.Printf("%s is listening on port %d\n", fn, port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		fmt.Println("Error starting file server:", err)
	}
}

func secretFilterHandler(next http.HandlerFunc) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		getSecret := r.URL.Query().Get(QueryParamKey)
		if getSecret != secret {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}

func fileListHandler(w http.ResponseWriter, r *http.Request) {
	files, err := os.ReadDir(path)
	if err != nil {
		http.Error(w, "Unable to list files", http.StatusInternalServerError)
		return
	}

	result := make([]string, 0)
	w.Header().Set("Content-Type", "application/json")
	for _, f := range files {
		if !f.IsDir() {
			result = append(result, f.Name())
		}
	}

	marshal, err := json.Marshal(map[string]interface{}{
		"data": result,
	})
	if err != nil {
		http.Error(w, "Unable to list files", http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(marshal)
}

func fileDownloadHandler(w http.ResponseWriter, r *http.Request) {
	fileName := strings.TrimPrefix(r.URL.Path, PathDownload)
	filePath := filepath.Join(path, fileName)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	http.ServeFile(w, r, filePath)
}
