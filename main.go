package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"net/http"
)

func main() {
	var (
		addr = flag.String("addr", ":8888", "HTTP address to bind to")
	)
	flag.Parse()

	http.DefaultServeMux.HandleFunc("/", rootHandler)
	http.DefaultServeMux.HandleFunc("/upload", uploadHandler)

	http.ListenAndServe(*addr, nil)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	if err := rootTmpl.Execute(w, struct {
		Message string
	}{
		Message: "",
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(2 << 20) // max 2 MB

	formfile, header, err := r.FormFile("csvfile")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer formfile.Close()

	file, err := header.Open()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var res = struct {
		Name string `json:"name"`
		Size int64  `json:"size"`
		Body string `json:"body"`
	}{
		Name: header.Filename,
		Size: header.Size,
		Body: string(data),
	}
	w.WriteHeader(http.StatusOK)
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
