package main

import (
	"context"
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Manager struct {
	template *template.Template
	dynamodbClient *dynamodb.Client
	tableName string
}

type Data struct {
	userId int
	userAge int
	userName string
}

func newManager() *Manager {
	awsConfig, awsConfigErr := config.LoadDefaultConfig(context.TODO())
	if awsConfigErr != nil {
		log.Println("AWS CONFIG ERROR")
	}
	return &Manager{
		template: template.Must(template.ParseGlob("html/*.html")),
		dynamodbClient: dynamodb.NewFromConfig(awsConfig),
		tableName: "go-practice",
	}
}

func (manager *Manager) run() {
	http.HandleFunc("/", manager.rootHandler)
	http.HandleFunc("/get", manager.getHandler)
	http.HandleFunc("/insert", manager.insertHandler)
	http.HandleFunc("/delete", manager.deleteHandler)
	http.HandleFunc("/main.js", handleStatic("main.js"))
	http.ListenAndServe(":8000", nil)
}

func handleStatic(fileName string) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		http.ServeFile(writer, request, fileName)
	}
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

func (manager *Manager) getHandler(writer http.ResponseWriter, request *http.Request) {
	response, getItemErr := manager.dynamodbClient.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(manager.tableName),
		Key: map[string]types.AttributeValue{
			"userId": &types.AttributeValueMemberN{Value: "1"},
		},
	})

	if getItemErr != nil {
		log.Println("GET ITEM FAILED")
		writer.WriteHeader(400);
		return
	}

	if response == nil || response.Item == nil {
		log.Println("NO REPONSE OR NO SUCH ITEM")
		writer.WriteHeader(500);
		return
	}

	data := &Data{}
	unmarshalErr := attributevalue.UnmarshalMap(response.Item, data)
	if unmarshalErr != nil {
		log.Println("GET ITEM UNMARSHAL ERROR")
		writer.WriteHeader(500);
		return
	}

	jsonData, jsonErr := json.Marshal(data)
	if jsonErr != nil {
		log.Println("JSON MARSHAL ERROR")
		writer.WriteHeader(500);
		return
	}

	writer.Write(jsonData)
}

func (manager *Manager) insertHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte("insertHandler was called"))
}

func (manager *Manager) deleteHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte("deleteHandler was called"))
}