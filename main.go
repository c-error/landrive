package main

/*
#include "embed/ansi.c"
*/
import "C"
import (
	"fmt"
	// "io"
	"time"
	// "unsafe"
	"net/http"
	"os"
	// "path/filepath"
	// "path"
	// "compress/gzip"
	"strings"
	// "syscall"
	// "mime"
	"encoding/base64"
	_ "embed"
)

var (

	// IP = "__DEMO__:8080"
	DIR = "./"
	port  = "8080"
	DECODE = strings.NewReplacer(
		"<0>", "#",
		"<1>", "%",
		"<2>", "&",
		"<3>", "+",
		"<4>", ";",
	)
	ENCODE = strings.NewReplacer(
		"#", "<0>",
		"%", "<1>",
		"&", "<2>",
		"+", "<3>",
		";", "<4>",
	)
	SLASH = strings.NewReplacer(
		"/", " /",
	)
	
	//go:embed embed/main.css
	CSS string
	//go:embed embed/main.js
	JS string
	//go:embed embed/login.js
	LJS string
	//go:embed embed/chat.js
	CJS string
	//go:embed embed/main.png
	PNG []byte
	icon = base64.StdEncoding.EncodeToString(PNG)

	//go:embed embed/font/CourierPrime-Regular.ttf
	FONT []byte
	font = base64.StdEncoding.EncodeToString(FONT)





	clint = NewStringSet()
)

const (
	// USER = "user"
	pass = "admin"
	ZIP = 100 * 1024 // 100 KB

	_DWN_ = `
	<!doctype html>
	<html lang="en">
	<head>
		<meta charset="utf-8" />
		<meta name="viewport" content="width=device-width,initial-scale=1" />
		<link rel="icon" href="data:image/png;base64,%s" />
		<title>LanDrive:/%s</title>
		<style>
			%s
		</style>
	</head>
	<body>
		<dwn>
			<info>
				<p>File Info: ...</p>
				<sub-info>
					<data><c>Name:</c><b>%s</b></data>
					<data><c>Type:</c><b>%s</b></data>
					<data><c>Size:</c><b>%s</b></data>
					<data><c>Date:</c><b>%s</b></data>
				</sub-info>
				<a href="%s">DOWNLOAD</a>
			</info>
		</dwn>
	</body>
	</html>
	`

	_ERR_ = `
	<!doctype html>
	<html lang="en">
	<head>
		<meta charset="utf-8" />
		<meta name="viewport" content="width=device-width,initial-scale=1" />
		<link rel="icon" href="data:image/png;base64,%s" />
		<title>LanDrive:/%s</title>
		<style>
			%s
		</style>
	</head>
	<body>
		%s
	</body>
	</html>
	`
)



// arr := []string{"Hello", "Hii", "Bye"}


// func func_exists(dir string) bool {

//     sdir := C.CString(dir)
//     defer C.free(unsafe.Pointer(sdir))
//     return C.is_exists(sdir) != 0

// }
















type StringSet struct {
	items map[string]struct{}
}

func NewStringSet() *StringSet {
	return &StringSet{
		items: make(map[string]struct{}),
	}
}

// Add a string (O(1) operation)
func (ss *StringSet) Add(s string) {
	ss.items[s] = struct{}{}
}

// Remove a string (O(1) operation)
func (ss *StringSet) Remove(s string) {
	delete(ss.items, s)
}

// Check if string exists (O(1) operation)
func (ss *StringSet) Contains(s string) bool {
	_, exists := ss.items[s]
	return exists
}

// Case-insensitive check (O(n) operation - unavoidable)
func (ss *StringSet) ContainsFold(s string) bool {
	s = strings.ToLower(s)
	for key := range ss.items {
		if strings.ToLower(key) == s {
			return true
		}
	}
	return false
}

// Get all items as slice
func (ss *StringSet) GetAll() []string {
	result := make([]string, 0, len(ss.items))
	for key := range ss.items {
		result = append(result, key)
	}
	return result
}

// Get length
func (ss *StringSet) Length() int {
	return len(ss.items)
}























func func_exists(path string) bool {

	info, err := os.Stat(path)
	if err != nil || !info.IsDir() {
		return false
	} else {
		return true
	}
}


// func func_decode(path string) string {

// 	// # -> <0>
// 	// % -> <1>
// 	// & -> <2>
// 	// + -> <3>
// 	// ; -> <4>

// 	temp := strings.ReplaceAll(path, "<0>", "#")
// 	temp = strings.ReplaceAll(temp, "<1>", "%")
// 	temp = strings.ReplaceAll(temp, "<2>", "&")
// 	temp = strings.ReplaceAll(temp, "<3>", "+")
// 	temp = strings.ReplaceAll(temp, "<4>", ";")
// 	return temp
// }

func func_decode(path string) string {
    return DECODE.Replace(path)
}

func func_encode(path string) string {
	return ENCODE.Replace(path)
}

// func func_log(a string, b string, c string) {
func func_log(color string, addr string, mode string, path string) {
	fmt.Printf("%s%s - %s - %s - %s \033[0m\n", color, time.Now().Format("15:04:05"), addr, mode, path)
}











func main() {

	C.enable_ansi()

	clint.Add("127.0.0.1")



	// Add items
	// set.Add("192.168.0.113")


	// Check existence
	// fmt.Println("Contains 'Apple':", set.Contains("Apple"))      // true
	// fmt.Println("Contains 'apple':", set.Contains("apple"))      // false
	// fmt.Println("ContainsFold 'apple':", set.ContainsFold("apple")) // true

	// Remove item
	// set.Remove("Banana")
	// fmt.Println("After removal:", set.GetAll())

	// Length
	// fmt.Println("Length:", set.Length())











	// Check if at least one argument is provided
	if len(os.Args) >= 3 { port = os.Args[1]; DIR = os.Args[2] }

	// fmt.Printf("IS_DIR: %s", func_exists(DIR))
	
	if func_exists(DIR) {

		http.HandleFunc("/path", handler_index)
		http.HandleFunc("/login", handler_login)
		http.HandleFunc("/set", handler_upload)
		http.HandleFunc("/get", handler_download)
		http.HandleFunc("/chat", handler_chat)
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

			w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate") // HTTP 1.1
			w.Header().Set("Pragma", "no-cache") // HTTP 1.0
			w.Header().Set("Expires", "0") // Proxies
			r.Write(w)
		})
		// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// 	http.Redirect(w, r, "/path?=/", http.StatusMovedPermanently)
		// })
		
		fmt.Printf("\nport : \033[92m%s\033[0m ...\ndir : \033[32m%s\033[0m ...\n\n", port, DIR)

		ip := fmt.Sprintf(":%s", port)
		if err := http.ListenAndServe(ip, nil); err != nil {
			
			fmt.Printf("\033[91mError @ %s !!\033[0m\n", port)
			os.Exit(0)
		}

	} else {

		fmt.Printf("\033[91mError @ %s !!\033[0m\n", DIR)
		os.Exit(0)
	}
}

































