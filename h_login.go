package main

import (
	"fmt"
	"net/http"
	"strings"
	// "syscall"
)

// const(

// )


func handler_login(w http.ResponseWriter, r *http.Request) {

	get_ip := r.RemoteAddr
	get_ip = get_ip[:strings.LastIndexByte(get_ip, ':')]

	if clint.Contains(get_ip) {

		func_log("\033[97m", r.RemoteAddr, "[MATCH]", get_ip)
		http.Redirect(w, r, "/path?fo=/", http.StatusFound)
		return
	}

	if r.URL.Query().Get("pin") == pass {

		clint.Add(get_ip)

		func_log("\033[97m", r.RemoteAddr, "[MATCH]", get_ip)
		http.Redirect(w, r, "/path?fo=/", http.StatusFound)
		return
		
	} else {

		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate") // HTTP 1.1
		w.Header().Set("Pragma", "no-cache") // HTTP 1.0
		w.Header().Set("Expires", "0") // Proxies
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		html := []byte(fmt.Sprintf(login_body, icon, font, CSS, icon, get_ip, LJS))
		w.Write(html)
		return
	}
}

