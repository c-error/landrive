package main


import (
	"fmt"
	// "io"
	// "time"
	// "unsafe"
	"net/http"
	"os"
	"path/filepath"
	"path"
	// "compress/gzip"
	// "strings"
	// "syscall"
	"mime"
	// "encoding/base64"
	// _ "embed"
)



func handler_download(w http.ResponseWriter, r *http.Request) {


	fmt.Println("CLINTL_IP:", r.RemoteAddr)

	mode_gt := r.URL.Query().Get("fi")

	if mode_gt != "" {
		build_path := filepath.Join(DIR, path.Clean(func_decode(mode_gt)))

		fileInfo, err := os.Stat(build_path)
		if err != nil || fileInfo.IsDir() {

			w.WriteHeader(http.StatusNotFound)
			html := fmt.Sprintf(
				_ERR_, 
				icon, 
				"/404", 
				CSS, 
				`<a>Fileinfo does not exist. </a><a class="error" href='/'>Return to home.</a>`,
			)
			w.Write([]byte(html))
			return
		}
		
		file, err := os.Open(build_path)
		if err != nil {

			w.WriteHeader(http.StatusNotFound)
			html := fmt.Sprintf(
				_ERR_, 
				icon, 
				"/404", 
				CSS, 
				`<a>File does not exist. </a><a class="error" href='/'>Return to home.</a>`,
			)
			w.Write([]byte(html))
			return
		}
		defer file.Close()
		
		contentType := mime.TypeByExtension(filepath.Ext(build_path))
		if contentType == "" {
			contentType = "application/octet-stream"
		}
		
		file_name := filepath.Base(build_path)
		mod_time := fileInfo.ModTime()

		w.Header().Set("Content-Type", contentType)
		w.Header().Set("Content-Length", string(fileInfo.Size()))
		w.Header().Set("Content-Disposition", "inline; filename=\""+file_name+"\"")
		w.Header().Set("Last-Modified", mod_time.UTC().Format(http.TimeFormat))
		w.Header().Set("ETag", fmt.Sprintf("\"%x-%x\"", mod_time.UnixNano(), fileInfo.Size()))
		
		http.ServeContent(w, r, file_name, mod_time, file)
	} else {

		fmt.Fprint(w, "error ...")
	}
}