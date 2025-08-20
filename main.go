package main

/*
#include "embed/ansi.c"
*/
import "C"
import (
	"fmt"
	"time"
	"net/http"
	"net"
	"os"
	"strings"
	"encoding/base64"
	_ "embed"
)

var (
	pass = "admin"
	root = "./"
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
	
	ZIP = 100 * 1024 // 100 KB
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

	fmt.Printf("| %s%s - %s - %s - %s \033[0m\n", color, time.Now().Format("15:04:05"), addr, mode, path)
}

func main() {

	C.enable_ansi() // enable console ansi
	clint.Add("127.0.0.1") // bypass login on localhost

	// get and set arguments
	if len(os.Args) >= 4 {

		port = os.Args[1]
		pass = os.Args[2]
		root = os.Args[3]
	}

	// set handler
	if func_exists(root) {

		http.HandleFunc("/path", handler_index)
		http.HandleFunc("/login", handler_login)
		http.HandleFunc("/set", handler_upload)
		http.HandleFunc("/get", handler_download)
		http.HandleFunc("/chat", handler_chat)
		http.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {

			r.Write(w)
			return
		})
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			
			if r.URL.Path == "/" {
				http.Redirect(w, r, "/path?fo=/", http.StatusFound)
				return
			}
		})

		fmt.Println("\n+- LanDrive -------")
		fmt.Printf("| root: \033[94m%s\033[0m ...\n", root)
		fmt.Printf("| port: \033[94m%s\033[0m ...\n", port)
		fmt.Printf("| pass: \033[94m%s\033[0m ...\n", pass)
		fmt.Println("+------------------")
		fmt.Printf("| url: \033[92mhttp://127.0.0.1:%s/path?fo=/\033[0m ...\n", port)

		addrs, err := net.InterfaceAddrs()
		if err != nil {
			
			fmt.Printf("\033[91mError @ network interface !!\033[0m\n")
			os.Exit(1)
		}

		for _, addr := range addrs {
			// check for ip address [not loopback]
			if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
				if ip4 := ipNet.IP.To4(); ip4 != nil {
					
					fmt.Printf("| url: \033[92mhttp://%s:%s/path?fo=/\033[0m ...\n", ip4.String(), port) // print IPv4
				}
			}
		}
		fmt.Println("+- Logs -----------")

		ip := fmt.Sprintf(":%s", port)
		if err := http.ListenAndServe(ip, nil); err != nil {
			
			fmt.Printf("\033[91mError @ unknown port !!\033[0m\n")
			os.Exit(1)
		}

	} else {

		fmt.Printf("\033[91mError @ unknown root !!\033[0m\n", root)
		os.Exit(1)
	}
}


