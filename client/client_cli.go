package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

func main() {
	// Define a command-line flag for the file path
	filePath := flag.String("file", "", "Path to the file to upload (required)")
	flag.Parse()

	// Validate the file path
	if *filePath == "" {
		log.Fatal("Error: file path is required. Use the -file flag to specify the file path.")
	}

	// Open the file
	file, err := os.Open(*filePath)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer file.Close()

	// Create a buffer to store the request body
	var buf bytes.Buffer

	// Create a new multipart writer with the buffer
	w := multipart.NewWriter(&buf)

	// Create a new form field for the file
	fw, err := w.CreateFormFile("file", file.Name())
	if err != nil {
		log.Fatalf("Error creating form file: %v", err)
	}

	// Copy the contents of the file to the form field
	if _, err := io.Copy(fw, file); err != nil {
		log.Fatalf("Error copying file contents: %v", err)
	}

	// Close the multipart writer to finalize the request
	if err := w.Close(); err != nil {
		log.Fatalf("Error closing multipart writer: %v", err)
	}

	// Create the HTTP request
	req, err := http.NewRequest("POST", "http://localhost:8080/upload_multipart", &buf)
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}
	defer resp.Body.Close()

	// Print the response
	fmt.Printf("Response status: %s\n", resp.Status)
}
