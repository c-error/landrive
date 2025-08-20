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

const (
body = `
<!doctype html>
<html lang="en">
<head>
	<meta charset="utf-8" />
	<meta name="viewport" content="width=device-width,initial-scale=1" />
	<link rel="icon" href="data:image/png;base64,%s" />
	<title>LanDrive:/%s</title>
	<style>
		@font-face { 
			font-family: "Courier Prime";
			src: url(data:application/octet-stream;base64,%s);
		}
		%s
	</style>
	<script>
		const SVR_URL = "%s";
	</script>
</head>
<body>


	<index>
		<panel>
			<form id="uploadForm" enctype="multipart/form-data">
				<input type="file" name="files" id="fileInput" multiple>
				<button type="button" onclick="uploadFiles()">UPLOAD</button>
			</form>
			<gap></gap>
			<mid-bar>
				<search>
					<input type="text" id="filter_name" placeholder="Filter by Name ...">
					<p>|</p>
					<input type="text" id="filter_size" placeholder="Filter by Size ...">
				</search>
				<filter>
					<div>
						<input type="checkbox" id="filter_all">
						<a>:All</a>
						<input type="checkbox" id="filter_fo">
						<a>:Folder</a>
						<input type="checkbox" id="filter_fi">
						<a>:File</a>
					</div>
					<div>
						<button onclick="search_clr();">CLEAR</button>
						<p>|</p>
						<button onclick="svr_search();">SEARCH</button>
					</div>
				</filter>
			</mid-bar>
			<gap></gap>
			<loc>
				<a href='/path?fo=/'>HOME</a>
				<b>|</b>
				<a href='%s'>BACK</a>
				<b>|</b>
				<p>%s</p>
			</loc>
		</panel>

		<uploding id="uploding_shell" style="display: none;">
			<h1>Uploding: ...</h1>
			<upload-subshell id="uploding"></upload-subshell>
		</uploding>

		%s

	</index>

	<div id="popup"></div>
	<script>
		%s
	</script>
</body>
</html>
`
)

func handler_index(w http.ResponseWriter, r *http.Request) {

	get_ip := r.RemoteAddr
	get_ip = get_ip[:strings.LastIndexByte(get_ip, ':')]

	fmt.Println("CLINTL_IP:", r.RemoteAddr)

	if !clint.Contains(get_ip) {
		http.Redirect(w, r, "/login", http.StatusMovedPermanently)
		return
	}

	mode_fo := r.URL.Query().Get("fo")
	mode_fi := r.URL.Query().Get("fi")

	if mode_fo != "" {

		clean_url := path.Clean(mode_fo)
		build_path := filepath.Join(DIR, func_decode(clean_url))
		
		fmt.Printf("\033[91mDEBUG: %s !!\033[0m\n", clean_url)

		if func_exists(build_path) {
			
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
					_ERR_, 
					icon, 
					"/500", 
					CSS, 
					`<a>Folder path error. </a><a class="error" href='/'>Return to home.</a>`,
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
					_ERR_, 
					icon, 
					"/500", 
					CSS, 
					`<a>Internal server error. </a><a class="error" href='/'>Return to home.</a>`,
				)
				w.Write([]byte(html))
				return
			}
		
			dir := C.GoString(c_dir_data)
			C.free(unsafe.Pointer(c_dir_data))
		
			if dir[:3] == "_E_" {
				
				cmsg := fmt.Sprintf(`<a>%s. </a><a class="error" href='/'>Return to home.</a>`, dir[3:])
				html := fmt.Sprintf(
					_ERR_, 
					icon, 
					"/500", 
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

			// top := fmt.Sprintf(_TOP_, parent, SLASH.Replace(func_decode(clean_url)))
			html := []byte(fmt.Sprintf(
				body,
				icon,
				clean_url,
				font,
				CSS,
				clean_url,
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

			html := fmt.Sprintf(
				_ERR_, 
				icon, 
				"/404", 
				CSS, 
				`<a>Folder does not exist. </a><a class="error" href='/'>Return to home.</a>`,
			)
			w.Write([]byte(html))
			return
		}


	} else if mode_fi != "" {

		clean_url := path.Clean(mode_fi)
		build_path := filepath.Join(DIR, func_decode(clean_url))

		fileInfo, err := os.Stat(build_path)
		if err != nil || fileInfo.IsDir() {

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
			_DWN_,
			icon,
			"/"+file_name,
			CSS,
			file_name,
			filepath.Ext(build_path),
			file_size,
			mod_time.Format("15:04:05-02:01:2006"),
			download,
		)
		w.Write([]byte(html))

	} else {

		http.Redirect(w, r, "/path?fo=/", http.StatusMovedPermanently)
	}
}
