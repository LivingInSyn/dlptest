package main

import (
	"html/template"
	"io"
	"log"
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

// Serve the template
func serveTemplate(w http.ResponseWriter, r *http.Request) {
	// Sample list of downloadable files
	files := []File{
		{Name: "file1.txt"},
		{Name: "file2.txt"},
		{Name: "file3.txt"},
	}

	data := PageData{
		Dlfs:         files,
		UseSlack:     true,  // Set this based on your application logic
		SlackWebhook: "your-slack-webhook-url",  // Replace with actual webhook URL
	}

	tmpl, err := template.ParseFiles("layout.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
		file, _, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Error retrieving the file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Save the uploaded file
		dst, err := os.Create(filepath.Join("uploads", "uploaded-file.txt"))
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

	// Set the appropriate content type and trigger the file download
	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	w.Header().Set("Content-Type", "application/octet-stream")
	http.ServeFile(w, r, filepath.Join("uploads", fileName))
}

func main() {
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
