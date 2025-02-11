package main

import (
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// DLFile represents a file with a name and SHA1 hash.
type DLFile struct {
	Name string
	Hash string
}

// hashFile calculates the SHA1 hash of a single file.
func hashFile(filePath string) (string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("opening file: %w", err)
	}
	defer f.Close()

	h := sha1.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", fmt.Errorf("hashing file: %w", err)
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

// hashDirectory calculates and returns the SHA1 hash of every file in a directory.
func hashDirectory(dirPath string) (map[string]DLFile, error) {
	dlfiles := make(map[string]DLFile)
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("accessing %s: %w", path, err)
		}

		// Skip directories and other non-regular files
		if !info.Mode().IsRegular() {
			return nil
		}

		hash, err := hashFile(path)
		if err != nil {
			return fmt.Errorf("hashing %s: %w", path, err)
		}

		dlfile := DLFile{
			Name: info.Name(),
			Hash: hash,
		}
		dlfiles[info.Name()] = dlfile

		fmt.Printf("%s: %s\n", path, hash)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("walking the directory: %w", err)
	}

	return dlfiles, nil
}

func main() {
	dirPath := "path/to/your/directory"
	dlfiles, err := hashDirectory(dirPath)
	if err != nil {
		fmt.Printf("Error hashing directory: %v\n", err)
		return
	}

	// Output the result
	for name, dlfile := range dlfiles {
		fmt.Printf("File: %s, Hash: %s\n", name, dlfile.Hash)
	}
}
