package main

import (
	"context"
	"html/template"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type Manager struct {
	template *template.Template
	dynamodbClient *dynamodb.Client
}

func newManager() *Manager {
	awsConfig, awsConfigErr := config.LoadDefaultConfig(context.TODO())
	if awsConfigErr != nil {
		log.Println("AWS CONFIG ERROR")
	}
	return &Manager{
		template: template.Must(template.ParseGlob("html/*.html")),
		dynamodbClient: dynamodb.NewFromConfig(awsConfig),
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