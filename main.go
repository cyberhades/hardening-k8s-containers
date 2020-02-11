package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"text/template"
	"time"
)

type Page struct {
	Title string
	Body  []byte
}

type Note struct {
	Name    string
	Size    int64
	LastMod string
}

type NoteList struct {
	Notes     []Note
	DebugInfo string
}

type NoteContent struct {
	Name    string
	Content string
}

func addCookie(w http.ResponseWriter, name string, value string) {
	expire := time.Now().AddDate(0, 0, 1)
	cookie := http.Cookie{
		Name:    name,
		Value:   value,
		Expires: expire,
	}
	http.SetCookie(w, &cookie)
}

func loadPage(title string) (*Page, error) {
	filename := title
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func createNote(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {

		cookie := getPathCookie(r)

		note := []byte(r.FormValue("note"))
		err := ioutil.WriteFile(cookie.Value+"/"+r.FormValue("title"), note, 0644)
		if err != nil {
			panic(err)
		}
	}

	page, err := loadPage("pages/create.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	addCookie(w, "path", "notes")
	w.Write(page.Body)
}

func getPathCookie(r *http.Request) *http.Cookie {
	cookie, err := r.Cookie("path")
	if err != nil {
		cookie = &http.Cookie{}
	}
	return cookie
}

func listNotes(w http.ResponseWriter, r *http.Request) {

	cookie := getPathCookie(r)

	files, err := ioutil.ReadDir(cookie.Value)
	if err != nil {
		log.Fatal(err)
	}

	noteList := NoteList{}

	debugInfo := r.URL.Query().Get("debug")
	if len(debugInfo) > 0 {
		cmd := exec.Command("/bin/sh", "-c", "cd "+cookie.Value+";ls "+debugInfo)
		out, err := cmd.CombinedOutput()
		if err != nil {
			log.Print("cmd.Run() failed with %s\n", err)
		}
		noteList.DebugInfo = string(out)
	}

	var notes []Note
	for _, f := range files {
		notes = append(notes, Note{
			f.Name(),
			f.Size(),
			f.ModTime().Format("2006.01.02 15:04:05")})
	}

	noteList.Notes = notes

	addCookie(w, "path", "notes")
	t := template.New("list.html")
	t, err = t.ParseFiles("pages/list.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = t.Execute(w, noteList)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func viewNote(w http.ResponseWriter, r *http.Request) {

	cookie := getPathCookie(r)

	fileName := strings.TrimPrefix(r.URL.Path, "/view/")
	content, err := ioutil.ReadFile(cookie.Value + "/" + fileName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	noteContent := NoteContent{fileName, string(content)}

	//addCookie(w, "path", "notes")
	t := template.New("view.html")
	t, err = t.ParseFiles("pages/view.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = t.Execute(w, noteContent)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {

	page, err := loadPage("pages/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	addCookie(w, "path", "notes")
	w.Write(page.Body)
}

func main() {
	fs := http.FileServer(http.Dir("assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))
	http.HandleFunc("/create", createNote)
	http.HandleFunc("/list", listNotes)
	http.HandleFunc("/view/", viewNote)
	http.HandleFunc("/", indexHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
