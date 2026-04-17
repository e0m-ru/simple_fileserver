// Package server handles HTTP handlers and server logic
package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/e0m-ru/fileserver/internal/config"
)

func CollectHandlers() (*http.ServeMux, error) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/api/upload", uploadHandler)
	mux.HandleFunc("/api/files", filesHandler)
	mux.HandleFunc("/api/delete", deleteHandler)
	uploads := http.FileServer(http.Dir(config.Config.Os.Uploads))
	mux.Handle("/uploads/", http.StripPrefix("/uploads/", uploads))
	staticHandler := http.FileServer(http.FS(config.StaticFiles))
	mux.Handle("/static/", staticHandler)
	return mux, nil
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	err := config.Config.Os.TMPL.Execute(w, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(32 << 20) // 32MB
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		log.Print(err)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to get file", http.StatusBadRequest)
		log.Print(err)
		return
	}
	defer file.Close()

	// Create uploads directory if not exists
	os.MkdirAll(config.Config.Os.Uploads, 0755)
	// Save file
	dst, err := os.Create(filepath.Join(config.Config.Os.Uploads, handler.Filename))
	if err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		log.Print(err)
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		log.Print(err)
		return
	}

	fmt.Fprintf(w, "File %s uploaded successfully!", handler.Filename)
}

func filesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	files, err := os.ReadDir(config.Config.Os.Uploads)
	if err != nil {
		if os.IsNotExist(err) {
			json.NewEncoder(w).Encode([]string{})
			return
		}
		http.Error(w, "Failed to read uploads directory", http.StatusInternalServerError)
		log.Print(err)
		return
	}

	var fileList []map[string]interface{}
	for _, file := range files {
		if !file.IsDir() {
			info, err := file.Info()
			if err != nil {
				continue
			}
			fileList = append(fileList, map[string]interface{}{
				"name": file.Name(),
				"size": info.Size(),
			})
		}
	}

	json.NewEncoder(w).Encode(fileList)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}
	filename := r.FormValue("filename")
	if filename == "" {
		http.Error(w, "Filename is required", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	// Prevent directory traversal attacks
	if strings.Contains(filename, "..") || strings.Contains(filename, "/") || strings.Contains(filename, "\\") {
		http.Error(w, "Invalid filename", http.StatusBadRequest)
		return
	}

	filePath := filepath.Join(config.Config.Os.Uploads, filename)
	err = os.Remove(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			http.Error(w, "File not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to delete file", http.StatusInternalServerError)
		}
		return
	}

	fmt.Fprintf(w, `{"success": true, "message": "File deleted successfully"}`)
}
