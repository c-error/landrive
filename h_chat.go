package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

var CLINT_CHAT []string
var CHAT_COUNT = 0

func handler_chat(w http.ResponseWriter, r *http.Request) {

	get_ip := r.RemoteAddr
	get_ip = get_ip[:strings.LastIndexByte(get_ip, ':')]

	if !clint.Contains(get_ip) {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	get_req := r.URL.Query().Get("req")

	if get_req == "sync" {

		if CHAT_COUNT != 0 {

			msg_count, err := strconv.Atoi(r.URL.Query().Get("no"))
			if err != nil {
				fmt.Println("Conversion error:")
			}

			if CHAT_COUNT >= msg_count {

				w.Write([]byte(CLINT_CHAT[msg_count-1]))
				return

			} else {

				w.Write([]byte("x"))
				return
			}

		} else {

			w.Write([]byte("x"))
			return
		}

	} else if get_req == "add" {
		
		get_data := r.URL.Query().Get("data")
		get_name := r.URL.Query().Get("name")
			
		if get_data != "" && get_name != "" {

			build_log := fmt.Sprintf("%s @ %s", get_name, get_data)
			func_log("\033[97m", r.RemoteAddr, "[CHAT]", func_decode(build_log))

			build_chat := fmt.Sprintf(`<div><p>%s</p><a>> %s</a></div>`, get_name, get_data)

			CLINT_CHAT = append(CLINT_CHAT, build_chat)
			CHAT_COUNT = len(CLINT_CHAT)

			// send := []byte(build_chat)
			// w.Write("send")

			return
		}
	}

	html := []byte(fmt.Sprintf(chat_body, icon, font, CSS, CJS))
	w.Write(html)
	return
}

