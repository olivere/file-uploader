package main

import (
	"html/template"
)

var (
	rootTmpl = template.Must(template.New("root").Parse(`
		<!DOCTYPE html>
		<html>
		<head>
		  <meta charset="utf-8">
		  <title>Upload test</title>
		</head>
		<body>
		  <h1>Uploader</h1>
		  <form id="uploadForm">
			<input type="file" id="fileToUpload">
		  </form>
		  <pre><code id="output"></code></pre>
		  <script>
		  const input = document.getElementById('fileToUpload');
		  const output = document.getElementById('output');

		  input.addEventListener('change', function(e) {
			e.preventDefault();

			output.innerText = '';
			
			const data = new FormData();
			const file = e.target.files[0];
			data.append('csvfile', file, file.name);
		
			fetch('/upload', {
			  method: 'POST',
			  body: data
			}).then(
			  response => {
				console.log(response);
				return response.text().then(t => { output.innerText = t; });
			  }
			).catch(
			  error => console.log(error)
			);
		  });
		  </script>
		</body>
		</html>
		`))
)
