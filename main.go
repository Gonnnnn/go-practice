package main

import (
	"html/template"
	"net/http"
)

type Manager struct {
	template *template.Template
}

func newManager() *Manager {
	return &Manager{
		template: template.Must(template.ParseGlob("html/*.html")),
	}
}

func (manager *Manager) run() {
	http.HandleFunc("/", manager.rootHandler)
	http.HandleFunc("/select", manager.selectHandler)
	http.HandleFunc("/insert", manager.insertHandler)
	http.HandleFunc("/delete", manager.deleteHandler)
	http.ListenAndServe(":8000", nil)
}

func main(){
	newManager().run()
}

func (manager *Manager) rootHandler(writer http.ResponseWriter, request *http.Request) {
	executeError := manager.template.ExecuteTemplate(writer, "main.html", nil)
	if executeError != nil {
		writer.Write([]byte("error"))
	}
}

func (manager *Manager) selectHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte("selectHandler was called"))
}

func (manager *Manager) insertHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte("insertHandler was called"))
}

func (manager *Manager) deleteHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte("deleteHandler was called"))
}