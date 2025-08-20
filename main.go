package main

/*
#include "embed/ansi.c"
*/
import "C"
import (
	"fmt"
	"time"
	"net/http"
	"os"
	"strings"
	"encoding/base64"
	_ "embed"
)

var (

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

type StringSet struct {
	items map[string]struct{}
}

func NewStringSet() *StringSet {
	return &StringSet{
		items: make(map[string]struct{}),
	}
}

func (ss *StringSet) Add(s string) {
	ss.items[s] = struct{}{}
}

func (ss *StringSet) Remove(s string) {
	delete(ss.items, s)
}

func (ss *StringSet) Contains(s string) bool {
	_, exists := ss.items[s]
	return exists
}

func (ss *StringSet) ContainsFold(s string) bool {
	s = strings.ToLower(s)
	for key := range ss.items {
		if strings.ToLower(key) == s {
			return true
		}
	}
	return false
}

func (ss *StringSet) GetAll() []string {
	result := make([]string, 0, len(ss.items))
	for key := range ss.items {
		result = append(result, key)
	}
	return result
}

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

func func_decode(path string) string {
    return DECODE.Replace(path)
}

func func_encode(path string) string {
	return ENCODE.Replace(path)
}

func func_log(color string, addr string, mode string, path string) {
	fmt.Printf("%s%s - %s - %s - %s \033[0m\n", color, time.Now().Format("15:04:05"), addr, mode, path)
}

func main() {

	C.enable_ansi()

	clint.Add("127.0.0.1")

	if len(os.Args) >= 3 { port = os.Args[1]; DIR = os.Args[2] }

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
