package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"text/template"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
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
type DownloadTemplate struct {
	Dlfs         []DLFile
	UseS3        bool
	UseSlack     bool
	SlackWebhook string
}
type S3Config struct {
	BucketName string
	Region     string
	AccessKey  string
	SecretKey  string
	Expiration time.Duration
}

// PreSignedURLResponse represents the S3 pre-signed URL response
type PreSignedURLResponse struct {
	URL    string            `json:"url"`
	Fields map[string]string `json:"fields"`
}

var DLFiles map[string]DLFile
var s3Config = S3Config{
	BucketName: "your-bucket-name",
	Region:     "us-east-1",
	AccessKey:  "your-access-key",
	SecretKey:  "your-secret-key",
	Expiration: 15 * time.Minute,
}
var useS3 = false

func main() {
	// setup logger
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	// populate the test files
	populateDLFiles()
	//handle s3
	useS3 = configureS3()
	// Create upload directory if it doesn't exist
	if _, err := os.Stat(uploadPath); os.IsNotExist(err) {
		os.MkdirAll(uploadPath, os.ModePerm)
	}

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", serveTemplate)

	http.HandleFunc("/upload", uploadFile)
	http.Handle("/download/", http.StripPrefix("/download/", http.FileServer(http.Dir(downloadPath))))

	http.HandleFunc("/availableFiles", getAvailableFiles)
	http.HandleFunc("/generateS3Token", generateS3PreSignedURL)

	fmt.Println("Server listening on port 8080")
	log.Info().Err(http.ListenAndServe(":8080", nil)).Msg("http server exiting")
}

func configureS3() bool {
	s3configured := true
	s3region, exists := os.LookupEnv("S3REGION")
	if exists {
		s3Config.Region = s3region
	} else {
		s3configured = false
	}
	s3bucket, exists := os.LookupEnv(("S3BUCKET"))
	if exists {
		s3Config.BucketName = s3bucket
	} else {
		s3configured = false
	}
	//	AccessKey:  "your-access-key",
	s3AccessKey, exists := os.LookupEnv(("S3KEYID"))
	if exists {
		s3Config.AccessKey = s3AccessKey
	} else {
		s3configured = false
	}
	//SecretKey:  "your-secret-key",
	s3secret, exists := os.LookupEnv(("S3SECRET"))
	if exists {
		s3Config.SecretKey = s3secret
	} else {
		s3configured = false
	}
	return s3configured
}

func populateDLFiles() {
	dirfiles, err := hashDirectory(downloadPath)
	if err != nil {
		log.Fatal().Err(err).Msg("Couldn't hash the directory")
	}
	DLFiles = dirfiles
}

func serveTemplate(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("./templates/layout.html"))
	// create the data for the page
	dlfs := make([]DLFile, 0, len(DLFiles))
	for _, v := range DLFiles {
		dlfs = append(dlfs, v)
	}
	dt := DownloadTemplate{
		Dlfs:         dlfs,
		SlackWebhook: "",
		UseSlack:     false,
		UseS3:        useS3,
	}
	// execute the template
	s3Url := fmt.Sprintf("https://%s.s3.us-west-2.amazonaws.com/", s3Config.BucketName)
	w.Header().Set("Access-Control-Allow-Origin", s3Url)
	err := tmpl.Execute(w, dt)
	if err != nil {
		addLogData(r, log.Error().Err(err), "error executing template")
		log.Print(err.Error())
		http.Error(w, http.StatusText(500), 500)
	}
}

func getAvailableFiles(w http.ResponseWriter, r *http.Request) {
	addLogData(r, log.Info(), "returning files available for download")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(DLFiles)
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
	addLogData(r, log.Info(), "New Upload")
	//reqLogger(r)
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		addLogData(r, log.Warn(), "got a non-post to upload")
		return
	}

	// Parse multipart form, setting max memory for file uploads
	r.ParseMultipartForm(maxUploadSize)
	if r.ContentLength > maxUploadSize {
		addLogData(r, log.Error(), "File too big")
		http.Error(w, "File too large", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		addLogData(r, log.Error(), "Error retrieving the file")
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		fmt.Println(err)
		return
	}
	defer file.Close()

	// Sanitize filename to prevent path traversal attacks
	filename := filepath.Base(handler.Filename) // Extract only the file name
	if filename == "." || filename == "/" {
		addLogData(r, log.Error().Str("filename", filename), "Invalid Filename")
		http.Error(w, "Invalid filename", http.StatusBadRequest)
		return
	}

	uploadFile, err := os.OpenFile(filepath.Join(uploadPath, filename), os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		addLogData(r, log.Error().Err(err), "couldn't create output file")
		http.Error(w, "Error creating the file for upload", http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
	defer uploadFile.Close()
	defer os.Remove(filepath.Join(uploadPath, filename))

	// Copy the uploaded file to the server's filesystem
	_, err = io.Copy(uploadFile, file)
	if err != nil {
		addLogData(r, log.Error().Err(err), "Error copying the file")
		http.Error(w, "Error copying the file", http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	// check the hash
	hashval, err := hashFile(filepath.Join(uploadPath, filename))
	if err != nil {
		addLogData(r, log.Error().Err(err), "Error hashing uploaded file")
		// log.Printf("Error hashing uploaded file: %s", err)
		http.Error(w, "Error retrieving the file", http.StatusInternalServerError)
		return
	}
	dlfilekey := handler.Filename
	if hashval != DLFiles[dlfilekey].Hash {
		le := log.Error().
			Str("Upload File Name", uploadFile.Name()).
			Str("Calculated hash val", hashval).
			Str("Expected Value", DLFiles[dlfilekey].Hash)
		addLogData(r, le, "Hash values don't match")
		// log.Printf("Hash for %s doesn't match. Got %s, expected %s\n", uploadFile.Name(), hashval, DLFiles[dlfilekey].Hash)
		http.Error(w, "Hash for file received by server doesn't match sample.", http.StatusUnauthorized)
		return
	} else {
		addLogData(r, log.Info().Str("Upload filename", uploadFile.Name()), "Hash matches")
		// log.Printf("Hash for %s matches!\n", uploadFile.Name())
	}

	addLogData(r, log.Info().Str("Upload filename", handler.Filename), "File uploaded successfully")
	// fmt.Fprintf(w, "File uploaded successfully: %s\n", handler.Filename)
}

func generateS3PreSignedURL(w http.ResponseWriter, r *http.Request) {
	// Create a context
	ctx := context.Background()
	addLogData(r, log.Info(), "starting s3 presigned URL request")

	// get the filename from the parameter
	// /generateS3Token?filename=foo
	query := r.URL.Query()
	reqFile, exists := query["filename"]
	if !exists || len(reqFile) == 0 {
		//fmt.Println("filters not present")
		addLogData(r, log.Error(), "no filename in s3 URL request")
		http.Error(w, "No Filename", http.StatusBadRequest)
		return
	}

	// Load AWS configuration
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(s3Config.AccessKey, s3Config.SecretKey, "")),
		config.WithRegion(s3Config.Region),
	)
	if err != nil {
		addLogData(r, log.Error().Err(err), "unable to load AWS configuration")
		http.Error(w, "Unable to load AWS configuration", http.StatusInternalServerError)
		return
	}

	// Create an S3 service client
	s3Client := s3.NewFromConfig(cfg)
	presignClient := s3.NewPresignClient(s3Client)

	// Generate a unique file name (e.g., based on the current timestamp)
	// TODO: get the file name under test
	fileName := fmt.Sprintf("%s-%d.txt", reqFile[0], time.Now().Unix())
	addLogData(r, log.Info().Str("filename", fileName), "filename for request not yet created")

	// Generate a pre-signed URL for the S3 upload
	req := &s3.PutObjectInput{
		Bucket: aws.String(s3Config.BucketName),
		Key:    aws.String(fileName),
	}

	presignedURL, err := presignClient.PresignPutObject(ctx, req, s3.WithPresignExpires(s3Config.Expiration))
	if err != nil {
		addLogData(r, log.Error().Err(err), "unable to create pre-signed URL")
		http.Error(w, "Unable to create pre-signed URL", http.StatusInternalServerError)
		return
	}

	// Respond with the pre-signed URL
	response := PreSignedURLResponse{
		URL:    presignedURL.URL,
		Fields: map[string]string{"key": fileName},
	}

	// Convert response to JSON
	w.Header().Set("Content-Type", "application/json")
	addLogData(r, log.Info().Str("filename", fileName), "Build and returned a S3 URL")
	// log.Printf("Build and returned a S3 URL for %s\n", fileName)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		addLogData(r, log.Error().Err(err), "Unable to encode response")
		http.Error(w, "Unable to encode response", http.StatusInternalServerError)
		return
	}
}
