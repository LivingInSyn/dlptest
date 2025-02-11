package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"text/template"
)

const (
	uploadPath    = "./uploads"
	downloadPath  = "./static/downloads"
	maxUploadSize = 350 * 1024 * 1024 // 350 MB
)

type DLFile struct {
	Name string
	Hash string
}

type DownloadTemplate struct {
	Dlfs         []DLFile
	UseSlack     bool
	SlackWebhook string
}

var DLFiles map[string]DLFile
var SlackWebhook string

func main() {
	value, exists := os.LookupEnv("DLPTEST_SLACK_HOOK")
	if exists {
		SlackWebhook = value
	} else {
		SlackWebhook = ""
	}
	populateDLFiles()

	if _, err := os.Stat(uploadPath); os.IsNotExist(err) {
		os.MkdirAll(uploadPath, os.ModePerm)
	}

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", serveTemplate)
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

func serveTemplate(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("./templates/layout.html"))
	dlfs := make([]DLFile, 0, len(DLFiles))
	for _, v := range DLFiles {
		dlfs = append(dlfs, v)
	}
	dt := DownloadTemplate{
		Dlfs:         dlfs,
		SlackWebhook: SlackWebhook,
		UseSlack:     SlackWebhook != "",
	}
	err := tmpl.Execute(w, dt)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, http.StatusText(500), 500)
	}
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

	r.ParseMultipartForm(maxUploadSize)
	if r.ContentLength > maxUploadSize {
		http.Error(w, "File too large", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		log.Println("Error retrieving the file:", err)
		return
	}
	defer file.Close()

	filename := filepath.Base(handler.Filename)
	if filename == "." || filename == "/" {
		http.Error(w, "Invalid filename", http.StatusBadRequest)
		return
	}

	uploadFile, err := os.OpenFile(filepath.Join(uploadPath, filename), os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		http.Error(w, "Error creating the file for upload", http.StatusInternalServerError)
		log.Println("Error creating file:", err)
		return
	}
	defer uploadFile.Close()
	defer os.Remove(filepath.Join(uploadPath, filename))

	_, err = io.Copy(uploadFile, file)
	if err != nil {
		http.Error(w, "Error copying the file", http.StatusInternalServerError)
		log.Println("Error copying file:", err)
		return
	}

	hashval, err := hashFile(filepath.Join(uploadPath, filename))
	if err != nil {
		log.Printf("Error hashing uploaded file: %s", err)
		http.Error(w, "Error retrieving the file", http.StatusInternalServerError)
		return
	}

	expectedFile, exists := DLFiles[filename]
	if !exists {
		log.Printf("Uploaded file %s does not match any expected files.", filename)
		http.Error(w, "Uploaded file does not match any expected files.", http.StatusUnauthorized)
		return
	}

	if hashval != expectedFile.Hash {
		log.Printf("File hash mismatch: Expected [%s], Received [%s]", expectedFile.Hash, hashval)
		http.Error(w, "Hash for file received by server doesn't match sample.", http.StatusUnauthorized)
		return
	}

	log.Printf("File %s uploaded successfully with matching hash!", filename)
	fmt.Fprintf(w, "File uploaded successfully: %s\n", filename)
}
