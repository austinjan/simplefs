package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"

	"path"
	"path/filepath"

	"github.com/pin/tftp"
)

var (
	port      string
	directory string
)

func init() {
	flag.StringVar(&port, "p", "8080", "Specify the port to listen on")
	flag.StringVar(&directory, "d", ".", "Specify the directory to serve files from")
}

func readHandler(filename string, rf io.ReaderFrom) error {
	file, err := os.Open(filepath.Join(directory, filename))
	if err != nil {
		return err
	}
	defer file.Close()
	// Obtain file information to get the file size
	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	rf.(tftp.OutgoingTransfer).SetSize(fileInfo.Size()) // Optional: informs the client about the file size
	_, err = rf.ReadFrom(file)
	return err
}

func writeHandler(filename string, wt io.WriterTo) error {
	file, err := os.Create(filepath.Join(directory, filename))
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = wt.WriteTo(file)
	return err
}

func main() {
	flag.Parse()

	http.HandleFunc("/list", listFilesHandler)
	http.HandleFunc("/upload", uploadFileHandler)
	http.Handle("/files/", http.StripPrefix("/files", http.FileServer(http.Dir(directory))))

	fmt.Printf("Serving %s on HTTP port: %s\n", directory, port)
	fmt.Println("Endpoints:")
	fmt.Println("- List all files: /list")
	fmt.Println("- Upload a file: /upload (POST)")
	fmt.Println("- Access a file: /files/{file_name}")

	// Start TFTP server in a goroutine
	go func() {
		s := tftp.NewServer(readHandler, writeHandler)
		err := s.ListenAndServe(":69") // TFTP standard port
		if err != nil {
			fmt.Printf("Failed to start TFTP server: %s\n", err)
		}
	}()

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
		return
	}
}

func listFilesHandler(w http.ResponseWriter, r *http.Request) {
	files, err := os.ReadDir(directory)
	if err != nil {
		http.Error(w, "Unable to read directory", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	for _, file := range files {
		if !file.IsDir() {
			fmt.Fprintf(w, "%s\n", file.Name())
		}
	}
}

func uploadFileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	r.ParseMultipartForm(10 << 20) // Limit upload size

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	filePath := filepath.Join(directory, handler.Filename)
	dst, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Error saving the file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, "Error writing the file", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "File uploaded successfully: %s\n", path.Base(dst.Name()))
}
