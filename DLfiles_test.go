package main

import (
	"crypto/sha1"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestHashFile(t *testing.T) {
	// Create temporary file
	tmpFile, err := os.CreateTemp("", "hashfile-test")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	content := []byte("hello world")
	if _, err := tmpFile.Write(content); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	// Compute expected hash
	h := sha1.New()
	h.Write(content)
	expectedHash := fmt.Sprintf("%x", h.Sum(nil))

	// Run hashFile
	gotHash, err := hashFile(tmpFile.Name())
	if err != nil {
		t.Errorf("hashFile failed: %v", err)
	}

	if gotHash != expectedHash {
		t.Errorf("Expected hash %s, got %s", expectedHash, gotHash)
	}
}

func TestHashDirectory(t *testing.T) {
	// Create a temporary directory
	tmpDir, err := os.MkdirTemp("", "hashdir-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a couple of files in the temporary directory
	file1Path := filepath.Join(tmpDir, "file1.txt")
	file2Path := filepath.Join(tmpDir, "file2.txt")

	content1 := []byte("file1 content")
	content2 := []byte("file2 content")

	if err := os.WriteFile(file1Path, content1, 0644); err != nil {
		t.Fatalf("Failed to write file1: %v", err)
	}
	if err := os.WriteFile(file2Path, content2, 0644); err != nil {
		t.Fatalf("Failed to write file2: %v", err)
	}

	// Calculate expected hashes
	h1 := sha1.Sum(content1)
	expectedHash1 := fmt.Sprintf("%x", h1[:])

	h2 := sha1.Sum(content2)
	expectedHash2 := fmt.Sprintf("%x", h2[:])

	// Run hashDirectory
	files, err := hashDirectory(tmpDir)
	if err != nil {
		t.Fatalf("hashDirectory failed: %v", err)
	}

	if len(files) != 2 {
		t.Errorf("Expected 2 files in map, got %d", len(files))
	}

	f1, exists := files["file1.txt"]
	if !exists {
		t.Errorf("file1.txt missing from returned map")
	} else if f1.Hash != expectedHash1 {
		t.Errorf("Expected hash for file1.txt to be %s, got %s", expectedHash1, f1.Hash)
	}

	f2, exists := files["file2.txt"]
	if !exists {
		t.Errorf("file2.txt missing from returned map")
	} else if f2.Hash != expectedHash2 {
		t.Errorf("Expected hash for file2.txt to be %s, got %s", expectedHash2, f2.Hash)
	}
}
