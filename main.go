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

// Manager
// template 관리하기 위해 template 객체를, dynamodb 요청 보내기 위해 dynamodb client를 가지고 있다.
// table name 여기다가 넣은 별다른 이유는 없음. 그냥 한 파일에 때려박으려다 보니..
type Manager struct {
	template *template.Template
	dynamodbClient *dynamodb.Client
	tableName string
}

// struct의 필드가 소문자로 시작했을때는 dynamodb에서 가져온 데이터를 unmarshal 해도 기본값이 들어갔다. 하지만 대문자로 시작하니 값이 제대로 출력됨을 확인했다.
type Data struct {
	UserId int
	UserAge int
	UserName string
}

func newManager() *Manager {
	// aws config를 가져와야한다. config를 가져오는 절차가 있다. 먼저 환경변수에서 가져오고, 그다음 /{HOME}/.aws/confidentials?인가에서 가져온다. 나는 후자에 등록해놨다.
	awsConfig, awsConfigErr := config.LoadDefaultConfig(context.TODO())
	if awsConfigErr != nil {
		log.Println("AWS CONFIG ERROR")
	}
	return &Manager{
		template: template.Must(template.ParseGlob("html/*.html")),
		dynamodbClient: dynamodb.NewFromConfig(awsConfig),
		tableName: "go-practice", // 자신의 dynamodb에 맞게 수정할 것
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
	// Data의 field가 소문자로 시작하면, unmarshal 결과 필드의 기본값만 남아있게 된다. 즉, unmarshal이 반영되지 않는다.
	// dynamodb의 attribute들도 물론 소문자로 시작하게 해놨다. 아마 별다른 네이밍 규칙이 있나보다. 자동으로 첫 문자를 대문자화하는 듯 하다.
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