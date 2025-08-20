package main

const(
	path_body = `
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
			<top-bar>
				<div class="logo">
					<img src="data:image/png;base64,%s">
					<h1>LanDrive</h1>
				</div>
				<div>
					<a href="/chat">CHAT</a>
					<p>|</p>
					<a href="/logout">LOGOUT</a>
				</div>
			</top-bar>
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
	chat_body = `
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
					<p>|</p>
					<a href="/path?fo=/">HOME</a>
				</div>
			</cell>
		</chat>
		<script>
			%s
		</script>
	</body>
	</html>
	`
	login_body = `
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
	dwn_body = `
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
	error_body = `
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
	</head>
	<body>
		<div class="error">
			%s
			<a href='/path?fo=/'>Return to home.</a>
		</div>
	</body>
	</html>
	`
)










