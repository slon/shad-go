package main

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/readFile", readFileHandler)
	mux.HandleFunc("/copy", copyHandler)
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func readFileHandler(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("file")
	data, _ := ioutil.ReadFile(filename)

	// Infer the Content-Type of the file.
	contentType := http.DetectContentType(data[:512])

	// Get the file size.
	fileSize := strconv.FormatInt(int64(len(data)), 10)

	// Send the headers.
	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Length", fileSize)

	_, _ = w.Write(data)
}

func copyHandler(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("file")
	f, _ := os.Open(filename)
	defer func() { _ = f.Close() }()

	// Infer the Content-Type of the file.
	filePrefix := make([]byte, 1024)
	_, _ = io.ReadAtLeast(f, filePrefix, 512)
	contentType := http.DetectContentType(filePrefix)

	// Get the file size.
	fstat, _ := f.Stat()
	fileSize := strconv.FormatInt(fstat.Size(), 10)

	// Send the headers.
	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Length", fileSize)

	_, _ = f.Seek(0, 0)
	_, _ = io.Copy(w, f)
}
