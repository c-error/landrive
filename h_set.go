package main



import (
	"fmt"
	"io"
	// "time"
	// "unsafe"
	"net/http"
	"os"
	"path/filepath"
	"path"
	// "compress/gzip"
	// "strings"
	// "syscall"
	// "mime"
	// "encoding/base64"
	// _ "embed"
)



func handler_upload(w http.ResponseWriter, r *http.Request) {

	mode_up := r.URL.Query().Get("fo")
	// fmt.Println(mode_up)

	if mode_up != "" {

		clean_url := path.Clean(mode_up)
		build_path := filepath.Join(DIR, func_decode(clean_url))

		fileInfo, err := os.Stat(build_path)
		if err != nil || !fileInfo.IsDir() {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Only parse the headers, not the content
		reader, err := r.MultipartReader()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Process each part as a stream
		for {

			part, err := reader.NextPart()
			if err == io.EOF {
				break
			}
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// Get filename from query parameter (more reliable than multipart filename)
			filename := r.URL.Query().Get("filename")
			if filename == "" {
				filename = part.FileName()
			}
			if filename == "" {
				continue
			}

			// fmt.Println(build_path)

			// Create destination file
			dstPath := filepath.Join(build_path, filepath.Base(filename))
			dst, err := os.Create(dstPath)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// Stream directly to disk
			if _, err := io.Copy(dst, part); err != nil {
				dst.Close()
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			dst.Close()
		}
		fmt.Fprint(w, "Files uploaded successfully")

	} else {

		fmt.Fprint(w, "error ...")
	}
}

