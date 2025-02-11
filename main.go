//main.go

package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

const (
	uploadPath    = "./uploads"
	downloadPath  = "./static/downloads"
	maxUploadSize = 350 * 1024 * 1024 // 350 MB
)

type DLFile struct {
	Name string `json:"name"`
	Hash string `json:"hash"`
}

type DownloadTemplate struct {
	Dlfs         []DLFile
	UseSlack     bool
	SlackWebhook string
}

var DLFiles map[string]DLFile
var SlackWebhook string

func main() {
	// Load Slack Webhook if set in the environment
	SlackWebhook, _ = os.LookupEnv("DLPTEST_SLACK_HOOK")

	// Ensure upload directory exists
	if _, err := os.Stat(uploadPath); os.IsNotExist(err) {
		os.MkdirAll(uploadPath, os.ModePerm)
	}

	// Populate expected file hashes
	populateDLFiles()

	// Setup Routes
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	http.HandleFunc("/", serveTemplate)
	http.HandleFunc("/upload", uploadFile)
	http.Handle("/download/", http.StripPrefix("/download/", http.FileServer(http.Dir(downloadPath))))
	http.HandleFunc("/availableFiles", getAvailableFiles)

	// Start Server
	fmt.Println("Server listening on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func populateDLFiles() {
	dirfiles, err := hashDirectory(downloadPath)
	if err != nil {
		log.Fatalf("Couldn't hash the directory: %v", err)
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

	if err := tmpl.Execute(w, dt); err != nil {
		http.Error(w, "Template rendering error", http.StatusInternalServerError)
		log.Printf("Template error: %v", err)
	}
}

func getAvailableFiles(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(DLFiles); err != nil {
		log.Printf("Error encoding JSON: %v", err)
		http.Error(w, "Failed to retrieve available files", http.StatusInternalServerError)
	}
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Limit file size before parsing
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)

	err := r.ParseMultipartForm(maxUploadSize)
	if err != nil {
		http.Error(w, "File too large", http.StatusBadRequest)
		log.Println("File upload error:", err)
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
	if filename == "." || filename == "/" || filename == "" {
		http.Error(w, "Invalid filename", http.StatusBadRequest)
		return
	}

	// Save file to uploads directory
	savePath := filepath.Join(uploadPath, filename)
	outFile, err := os.Create(savePath)
	if err != nil {
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		log.Println("Error saving file:", err)
		return
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, file)
	if err != nil {
		http.Error(w, "Error copying file", http.StatusInternalServerError)
		log.Println("Error copying file:", err)
		return
	}

	// Compute file hash
	hashval, err := hashFile(savePath)
	if err != nil {
		http.Error(w, "Error hashing file", http.StatusInternalServerError)
		log.Println("Error hashing file:", err)
		os.Remove(savePath) // Remove the uploaded file if there's an error
		return
	}

	// Verify against expected hash
	expectedFile, exists := DLFiles[filename]
	if !exists {
		log.Printf("File %s does not match any expected files.", filename)
		http.Error(w, "File does not match expected files.", http.StatusUnauthorized)
		os.Remove(savePath) // Remove the uploaded file if it doesn’t match
		return
	}

	if hashval != expectedFile.Hash {
		log.Printf("File hash mismatch for %s: Expected [%s], Received [%s]", filename, expectedFile.Hash, hashval)
		http.Error(w, "File hash mismatch.", http.StatusUnauthorized)
		os.Remove(savePath) // Remove the uploaded file if there's a hash mismatch
		return
	}

	log.Printf("File %s uploaded successfully with matching hash!", filename)
	fmt.Fprintf(w, "File uploaded successfully: %s\n", filename)
}

// hashFile calculates the SHA-256 hash of a file
func hashFile(filePath string) (string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

// hashDirectory computes SHA-256 hashes for files in a directory
func hashDirectory(dirPath string) (map[string]DLFile, error) {
	files := make(map[string]DLFile)

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("Error accessing file %s: %v", path, err)
			return nil // Skip problematic files
		}

		if !info.Mode().IsRegular() {
			return nil
		}

		hash, err := hashFile(path)
		if err != nil {
			log.Printf("Error hashing file %s: %v", path, err)
			return nil
		}

		files[info.Name()] = DLFile{Name: info.Name(), Hash: hash}
		log.Printf("Hashed: %s -> %s", info.Name(), hash)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error walking the directory: %w", err)
	}

	return files, nil
}
