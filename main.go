package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

const (
	uploadPath    = "./uploads"          // Directory to store uploaded files
	downloadPath  = "./static/downloads" // Directory to serve files from
	maxUploadSize = 350 * 1024 * 1024    // 350 MB
)

type DLFile struct {
	Name string
	Hash string
}

var DLFiles map[string]DLFile

func main() {
	populateDLFiles()
	// Create upload directory if it doesn't exist
	if _, err := os.Stat(uploadPath); os.IsNotExist(err) {
		os.MkdirAll(uploadPath, os.ModePerm)
	}

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	http.HandleFunc("/upload", uploadFile)
	http.Handle("/download/", http.StripPrefix("/download/", http.FileServer(http.Dir(downloadPath))))

	http.HandleFunc("/availableFiles", getAvailableFiles)

	fmt.Println("Server listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func populateDLFiles() {
	dirfiles, err := hashDirectory(downloadPath)
	if err != nil {
		log.Panic("couldn't hash the dir", err)
	}
	DLFiles = dirfiles
}

func getAvailableFiles(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(DLFiles)
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse multipart form, setting max memory for file uploads
	r.ParseMultipartForm(maxUploadSize)
	if r.ContentLength > maxUploadSize {
		http.Error(w, "File too large", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		fmt.Println(err)
		return
	}
	defer file.Close()

	// Sanitize filename to prevent path traversal attacks
	filename := filepath.Base(handler.Filename) // Extract only the file name
	if filename == "." || filename == "/" {
		http.Error(w, "Invalid filename", http.StatusBadRequest)
		return
	}

	uploadFile, err := os.OpenFile(filepath.Join(uploadPath, filename), os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		http.Error(w, "Error creating the file for upload", http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
	defer uploadFile.Close()
	defer os.Remove(uploadPath)

	// Copy the uploaded file to the server's filesystem
	_, err = io.Copy(uploadFile, file)
	if err != nil {
		http.Error(w, "Error copying the file", http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	// check the hash
	hashval, err := hashFile(uploadPath)
	if hashval != DLFiles[uploadFile.Name()].Hash {
		log.Printf("Hash for %s doesn't match. Got %s, expected %s\n", uploadFile.Name(), hashval, DLFiles[uploadFile.Name()].Hash)
	} else {
		log.Printf("Hash for %s matches!\n", uploadFile.Name())
	}

	fmt.Fprintf(w, "File uploaded successfully: %s\n", handler.Filename)
}
