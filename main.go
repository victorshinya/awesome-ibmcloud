package main

import (
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"text/template"

	"github.com/gorilla/mux"
	gfm "github.com/shurcooL/github_flavored_markdown"
)

const (
	readmePath = "./README.md"
	tmplPath   = "./public/tmpl.html"
	indexPath  = "./public/index.html"
)

type content struct {
	Body string
}

func generateHandler(w http.ResponseWriter, r *http.Request) {
	// Pull the latest Git version
	exec.Command("git checkout -f").Output()
	exec.Command("git pull").Output()

	// Generate index.html file
	input, _ := ioutil.ReadFile(readmePath)
	body := string(gfm.Markdown(input))
	c := &content{Body: body}
	t := template.Must(template.ParseFiles(tmplPath))
	f, _ := os.Create(indexPath)
	t.Execute(f, c)
	w.Write([]byte("Page generated"))
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/generate", generateHandler)
	fs := http.FileServer(http.Dir("public"))
	r.Handle("/", fs)

	http.ListenAndServe(":8080", r)
}
