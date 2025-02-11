package main

import (
	"html/template"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
)

// File struct to store file data
type File struct {
	Name string
}

// Data structure for passing to templates
type PageData struct {
	Dlfs        []File
	UseSlack    bool
	SlackWebhook string
}

// Ensure the uploads directory exists
func ensureUploadsDir() {
	err := os.MkdirAll("uploads", os.ModePerm)
	if err != nil {
		log.Fatal("Failed to create uploads directory: ", err)
	}
}

// Serve the template with dynamic file listing
func serveTemplate(w http.ResponseWriter, r *http.Request) {
	// Ensure uploads directory exists
	ensureUploadsDir()

	// Dynamically list files in the "uploads" directory
	var files []File
	err := filepath.Walk("uploads", func(path string, info os.FileInfo, err error) error {
		if err != nil || !info.Mode().IsRegular() {
			return err
		}
		files = append(files, File{Name: info.Name()})
		return nil
	})
	if err != nil {
		http.Error(w, "Error reading uploaded files", http.StatusInternalServerError)
		return
	}

	// Prepare data to pass to the template
	data := PageData{
		Dlfs:         files,
		UseSlack:     true,  
		SlackWebhook: "your-slack-webhook-url", 
	}

	// Parse and execute the template
	tmpl, err := template.ParseFiles("layout.html")
	if err != nil {
		http.Error(w, "Error parsing template: "+err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, data)
}

// Handle file upload
func handleFileUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Parse the form to handle the file upload
		err := r.ParseMultipartForm(10 << 20) // 10 MB limit
		if err != nil {
			http.Error(w, "Unable to parse the form", http.StatusBadRequest)
			return
		}

		// Retrieve the file from the form
		file, fileHeader, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Error retrieving the file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Ensure uploads directory exists
		ensureUploadsDir()

		// Create a destination file
		dst, err := os.Create(filepath.Join("uploads", fileHeader.Filename))
		if err != nil {
			http.Error(w, "Error saving the file", http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		// Copy the contents of the file to the destination
		_, err = io.Copy(dst, file)
		if err != nil {
			http.Error(w, "Error copying the file", http.StatusInternalServerError)
			return
		}

		// Respond to the user that the file was uploaded successfully
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// Handle file download
func handleFileDownload(w http.ResponseWriter, r *http.Request) {
	fileName := r.URL.Path[len("/download/"):]

	// Open the file for reading
	file, err := os.Open(filepath.Join("uploads", fileName))
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	defer file.Close()

	// Set the appropriate content type
	mimeType := mime.TypeByExtension(filepath.Ext(fileName))
	if mimeType != "" {
		w.Header().Set("Content-Type", mimeType)
	} else {
		w.Header().Set("Content-Type", "application/octet-stream")
	}

	// Set the Content-Disposition header to trigger the file download
	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	http.ServeFile(w, r, filepath.Join("uploads", fileName))
}

func main() {
	// Serve the template, handle uploads and downloads
	http.HandleFunc("/", serveTemplate)
	http.HandleFunc("/upload", handleFileUpload)
	http.HandleFunc("/download/", handleFileDownload)

	// Start the server
	log.Println("Server starting on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Server failed to start: ", err)
	}
}
