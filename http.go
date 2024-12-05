package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func runHttpServer(port int) {
	http.HandleFunc("/files", fileListHandler)
	http.HandleFunc("/download/", fileDownloadHandler)

	fmt.Printf("File server started at http://localhost:%s\n", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		fmt.Println("Error starting file server:", err)
	}
}

// 查看文件列表的处理器
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

// 文件下载的处理器
func fileDownloadHandler(w http.ResponseWriter, r *http.Request) {
	fileName := strings.TrimPrefix(r.URL.Path, "/download/")
	filePath := filepath.Join(path, fileName)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	http.ServeFile(w, r, filePath)
}
