package main

import (
	"fmt"
	// "io"
	// "time"
	// "unsafe"
	"net/http"
	"strconv"
	// "os"
	// "path/filepath"
	// "path"
	// "compress/gzip"
	// "strings"
	// "syscall"
	// "mime"
	// "encoding/base64"
	// _ "embed"
)

const(
	_CHAT_ = `
	<!doctype html>
	<html lang="en">
	<head>
		<meta charset="utf-8" />
		<meta name="viewport" content="width=device-width,initial-scale=1" />
		<link rel="icon" href="data:image/png;base64,%s" />
		<title>LanDrive://chat</title>
		<style>
			@font-face { 
				font-family: "Courier Prime";
				src: url(data:application/octet-stream;base64,%s);
			}
			%s
		</style>
	</head>
	<body>
		<chat>
			<msg>
				<sub-shell id="server_msg"></sub-shell>
				<input id="end_beacon">
			</msg>
			<cell>
				<input id="user_name" placeholder="Name ...">
				<p>|</p>
				<div>
					<input id="user_input" placeholder="Text ....">
					<button onclick="send_chat();">SEND</button>
				</div>
			</cell>
		</chat>
		<script>
		

		%s

		</script>
	</body>
	</html>
	`
)


var CLINT_CHAT []string
var CHAT_COUNT = 0

func handler_chat(w http.ResponseWriter, r *http.Request) {

	get_req := r.URL.Query().Get("req")

	if get_req == "sync" {

		if CHAT_COUNT != 0 {

			msg_count, err := strconv.Atoi(r.URL.Query().Get("no"))
			if err != nil {
				fmt.Println("Conversion error:")
			}

			if CHAT_COUNT >= msg_count {

				// if msg_count == 0 {
				// 	msg_count++
				// }
	
				// fmt.Printf("NOW: %d -> %s\n", msg_count, CLINT_CHAT[msg_count-1])
	



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


			

			build_chat := fmt.Sprintf(`<div><p>%s</p><a>> %s</a></div>`, get_name, get_data)

			CLINT_CHAT = append(CLINT_CHAT, build_chat)
			CHAT_COUNT = len(CLINT_CHAT)

			send := []byte(build_chat)
			w.Write(send)

			// fmt.Println("Updated slice:", CLINT_CHAT)
			return
		}
	}

	html := []byte(fmt.Sprintf(_CHAT_, icon, font, CSS, CJS))
	w.Write(html)
	return
}

