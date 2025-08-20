package main

import (
	"fmt"
	"net/http"
	"strings"
	// "syscall"
)

const(
	_LOGIN_ = `
	<!doctype html>
	<html lang="en">
	<head>
		<meta charset="utf-8" />
		<meta name="viewport" content="width=device-width,initial-scale=1" />
		<link rel="icon" href="data:image/png;base64,%s" />
		<title>LanDrive://login</title>
		<style>
			@font-face { 
				font-family: "Courier Prime";
				src: url(data:application/octet-stream;base64,%s);
			}
			%s
		</style>
		<script>
			const SVR_PAR = new URLSearchParams(window.location.search);
		</script>
	</head>
	<body>
		<login>
			<shell>
				<img src="data:image/png;base64,%s">
				<end>
					<h1>LanDrive @-%s</h1>
					<cell>
						<a>🔐 pin:</a>
						<input id="pin" type="password" placeholder="...">
					</cell>
				</end>
			</shell>
		</login>
		<script>
			%s
		</script>
	</body>
	</html>
	`
)


func handler_login(w http.ResponseWriter, r *http.Request) {

	get_ip := r.RemoteAddr
	get_ip = get_ip[:strings.LastIndexByte(get_ip, ':')]


	if clint.Contains(get_ip) {
		http.Redirect(w, r, "/path?fo=/", http.StatusMovedPermanently)
		return
	}

	get_pass := r.URL.Query().Get("pin")

	if get_pass == pass {

		clint.Add(get_ip)

		http.Redirect(w, r, "/path?fo=/", http.StatusMovedPermanently)
		return
	}

	html := []byte(fmt.Sprintf(_LOGIN_, icon, font, CSS, icon, get_ip, LJS))
	w.Write(html)
	return
}

