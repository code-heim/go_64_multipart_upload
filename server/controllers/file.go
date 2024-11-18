package controllers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func FileUpload(w http.ResponseWriter, r *http.Request) {
	// Limit the size of the incoming request body to prevent abuse (e.g., 100 MB)
	const maxUploadSize = 100 << 20 // 100 MB
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)

	// Parse the form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	// Get the filename
	filename := r.FormValue("file")
	if filename == "" {
		http.Error(w, "Filename query parameter is required", http.StatusBadRequest)
		return
	}

	// Define the destination path
	uploadDir := "./uploads"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		http.Error(w, "Failed to create upload directory: "+err.Error(), http.StatusInternalServerError)
		return
	}
	destPath := filepath.Join(uploadDir, filename)

	// Create the destination file
	destFile, err := os.Create(destPath)
	if err != nil {
		http.Error(w, "Failed to create destination file: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer destFile.Close()

	// Copy the raw body to the destination file
	if _, err := io.Copy(destFile, r.Body); err != nil {
		http.Error(w, "Failed to save file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with a success message
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "File uploaded successfully: %s\n", filename)
}

func FileUploadMultipart(w http.ResponseWriter, r *http.Request) {
	// Limit the size of the incoming request body to prevent abuse (e.g., 100 MB)
	const maxUploadSize = 100 << 20 // 100 MB
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)

	// Parse the multipart form
	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		http.Error(w, "Failed to parse multipart form: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Retrieve the file from the form data
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Unable to retrieve file from form: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	fmt.Printf("Uploaded File: %+v\n", fileHeader.Filename)
	fmt.Printf("File Size: %+v\n", fileHeader.Size)
	fmt.Printf("MIME Header: %+v\n", fileHeader.Header)

	// Define the destination path
	uploadDir := "./uploads"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		http.Error(w, "Failed to create upload directory: "+err.Error(), http.StatusInternalServerError)
		return
	}
	destPath := filepath.Join(uploadDir, fileHeader.Filename)

	// Create the destination file
	destFile, err := os.Create(destPath)
	if err != nil {
		http.Error(w, "Failed to create destination file: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer destFile.Close()

	// Copy the file's content to the destination file
	if _, err := io.Copy(destFile, file); err != nil {
		http.Error(w, "Failed to save file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with a success message
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "File uploaded successfully: %s\n", fileHeader.Filename)
}
