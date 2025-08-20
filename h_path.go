package main
/*
#include "embed/main.c"
*/
import "C"
import (
	"fmt"
	"unsafe"
	"net/http"
	"os"
	"path/filepath"
	"path"
	"compress/gzip"
	"strings"
	"syscall"
)

// const (

// )

func handler_index(w http.ResponseWriter, r *http.Request) {

	get_ip := r.RemoteAddr
	get_ip = get_ip[:strings.LastIndexByte(get_ip, ':')]

	// fmt.Println("CLINTL_IP:", r.RemoteAddr)

	if !clint.Contains(get_ip) {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	mode_fo := r.URL.Query().Get("fo")
	mode_fi := r.URL.Query().Get("fi")

	if mode_fo != "" {

		clean_url := path.Clean(mode_fo)
		decode_url := func_decode(clean_url)
		build_path := filepath.Join(root, decode_url)

		if func_exists(build_path) {
			
			func_log("\033[97m", r.RemoteAddr, "[PATH]", decode_url)

			FILTER_DATA := r.URL.Query().Get("f")

			FILTER_STATE := 0
			FILTER_TYPE := 0
			FILTER_NAME := ""
			FILTER_SIZE := ""

			if FILTER_DATA != "" {

				FILTER_STATE = 1
				filter := strings.Split(FILTER_DATA, ",")
				
				if filter[0] != "" {
					if filter[0] == "fo" {
						FILTER_TYPE = 1
					} else if filter[0] == "fi" {
						FILTER_TYPE = 2
					} else {
						FILTER_TYPE = 0
					}
				}
				if len(filter) >= 2 {
					FILTER_NAME = filter[1]
				}
				if len(filter) >= 3 {
					FILTER_SIZE = filter[2]
				}
			}

			C_BUILD_PATH, err := syscall.UTF16PtrFromString(build_path)
			if err != nil {
				html := fmt.Sprintf(
					error_body, 
					icon, 
					"/500", 
					font,
					CSS, 
					`<b>Folder path error. </b>`,
				)
				w.Write([]byte(html))
				return
			}

			C_CLINT_IP := C.CString(r.RemoteAddr)
			defer C.free(unsafe.Pointer(C_CLINT_IP))

			C_FILTER_STATE := C.int(FILTER_STATE)
			C_FILTER_TYPE := C.int(FILTER_TYPE)
			C_FILTER_NAME := C.CString(FILTER_NAME)
			C_FILTER_SIZE := C.CString(FILTER_SIZE)

			defer C.free(unsafe.Pointer(C_FILTER_NAME))
			defer C.free(unsafe.Pointer(C_FILTER_SIZE))

			c_dir_data := C.get_data_list(
				(*C.WCHAR)(C_BUILD_PATH),
				C_FILTER_STATE,  // filter 1->on / 0->off
				C_FILTER_TYPE,   // type 0->all / 1->fo / 2->fi
				C_FILTER_NAME,   // filter by name
				C_FILTER_SIZE,   // filter by size
				C_CLINT_IP,      // clint ip
			)

			if c_dir_data == nil {

				C.free(unsafe.Pointer(c_dir_data))
				html := fmt.Sprintf(
					error_body, 
					icon, 
					"/500",
					font,
					CSS, 
					`<b>Internal server error. </b>`,
				)
				w.Write([]byte(html))
				return
			}
		
			dir := C.GoString(c_dir_data)
			C.free(unsafe.Pointer(c_dir_data))
		
			if dir[:3] == "_E_" {
				
				cmsg := fmt.Sprintf(`<b>%s. </b>`, dir[3:])
				html := fmt.Sprintf(
					error_body, 
					icon, 
					"/500",
					font,
					CSS, 
					cmsg,
				)
				w.Write([]byte(html))
				return
			}

			// Set headers to prevent caching
			w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate") // HTTP 1.1
			w.Header().Set("Pragma", "no-cache") // HTTP 1.0
			w.Header().Set("Expires", "0") // Proxies
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			
			parent := clean_url[:strings.LastIndexByte(clean_url, '/')]
			if parent == "" { parent = "/" }
			parent = fmt.Sprintf("/path?fo=%s", parent)

			html := []byte(fmt.Sprintf(
				path_body,
				icon,
				clean_url,
				font,
				CSS,
				clean_url,
				icon,
				parent,
				SLASH.Replace(func_decode(clean_url)),
				dir,
				JS,
			))
		
			if len(html) >= ZIP && strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
				
				w.Header().Set("Content-Encoding", "gzip")
				gz, err := gzip.NewWriterLevel(w, gzip.BestCompression)
				if err != nil {

					w.Write(html)
					return
				}
				defer gz.Close()
				if _, err := gz.Write(html); err != nil {

					w.Write(html)
					return
				}
	
			} else {

				w.Write(html)
				return
			}
		} else {

			func_log("\033[91m", r.RemoteAddr, "[PATH]", decode_url)

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


	} else if mode_fi != "" {

		clean_url := path.Clean(mode_fi)
		decode_url := func_decode(clean_url)
		build_path := filepath.Join(root, decode_url)

		func_log("\033[97m", r.RemoteAddr, "[PATH]", decode_url)

		fileInfo, err := os.Stat(build_path)
		if err != nil || fileInfo.IsDir() {

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

		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate") // HTTP 1.1
		w.Header().Set("Pragma", "no-cache") // HTTP 1.0
		w.Header().Set("Expires", "0") // Proxies
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		get_data := C.get_file_size(C.size_t(fileInfo.Size()))
		file_size := C.GoString(get_data)
		C.free(unsafe.Pointer(get_data))

		file_name := filepath.Base(build_path)
		mod_time := fileInfo.ModTime()

		download := fmt.Sprintf("/get?fi=%s", clean_url)
		html := fmt.Sprintf(
			dwn_body,
			icon,
			"/"+file_name,
			font,
			CSS,
			file_name,
			filepath.Ext(build_path),
			file_size,
			mod_time.Format("15:04:05-02:01:2006"),
			download,
		)
		w.Write([]byte(html))

	} else {

		http.Redirect(w, r, "/path?fo=/", http.StatusFound)
	}
}
