package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"path"
	"mime"
)



func handler_download(w http.ResponseWriter, r *http.Request) {

	mode_gt := r.URL.Query().Get("fi")

	if mode_gt != "" {

		decode_url := func_decode(path.Clean(mode_gt))
		build_path := filepath.Join(root, decode_url)

		func_log("\033[97m", r.RemoteAddr, "[GET]", decode_url)

		fileInfo, err := os.Stat(build_path)
		if err != nil || fileInfo.IsDir() {

			w.WriteHeader(http.StatusNotFound)
			html := fmt.Sprintf(
				error_body, 
				icon, 
				"/404",
				font,
				CSS, 
				`<b>Fileinfo does not exist. </b>`,
			)
			w.Write([]byte(html))
			return
		}
		
		file, err := os.Open(build_path)
		if err != nil {

			w.WriteHeader(http.StatusNotFound)
			html := fmt.Sprintf(
				error_body, 
				icon, 
				"/404",
				font,
				CSS, 
				`<b>File does not exist. </b>`,
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

		func_log("\033[91m", r.RemoteAddr, "[GET]", func_decode(path.Clean(mode_gt)))

		w.WriteHeader(http.StatusNotFound)
		html := fmt.Sprintf(
			error_body, 
			icon, 
			"/404",
			font,
			CSS, 
			`<b>Folder does not exist. </b>`,
		)
		w.Write([]byte(html))
		return
	}
}