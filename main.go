package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// S3Config contains the S3 credentials and configuration
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

var s3Config = S3Config{
	BucketName: "your-bucket-name",
	Region:     "us-east-1",
	AccessKey:  "your-access-key",
	SecretKey:  "your-secret-key",
	Expiration: 15 * time.Minute,
}

// GenerateS3PreSignedURL generates a pre-signed URL for uploading files to S3
func GenerateS3PreSignedURL(w http.ResponseWriter, r *http.Request) {
	// Create a context
	ctx := context.Background()

	// Load AWS configuration
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(s3Config.AccessKey, s3Config.SecretKey, "")),
		config.WithRegion(s3Config.Region),
	)
	if err != nil {
		http.Error(w, "Unable to load AWS configuration", http.StatusInternalServerError)
		return
	}

	// Create an S3 service client
	s3Client := s3.NewFromConfig(cfg)
	presignClient := s3.NewPresignClient(s3Client)

	// Generate a unique file name (e.g., based on the current timestamp)
	fileName := fmt.Sprintf("test-file-%d.txt", time.Now().Unix())

	// Generate a pre-signed URL for the S3 upload
	req := &s3.PutObjectInput{
		Bucket: aws.String(s3Config.BucketName),
		Key:    aws.String(fileName),
	}

	presignedURL, err := presignClient.PresignPutObject(ctx, req, s3.WithPresignExpires(s3Config.Expiration))
	if err != nil {
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
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Unable to encode response", http.StatusInternalServerError)
		return
	}
}

// ServeLayout serves the layout HTML page
func ServeLayout(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/layout.html")
	if err != nil {
		http.Error(w, "Unable to load layout", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Unable to render layout", http.StatusInternalServerError)
	}
}

func main() {
	// Serve static files (CSS, JS, etc.) from the "static" folder
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Serve the layout HTML page
	http.HandleFunc("/", ServeLayout)

	// Handle pre-signed URL generation
	http.HandleFunc("/generate-s3-token", GenerateS3PreSignedURL)

	log.Println("Server starting on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Server failed to start: ", err)
	}
}
