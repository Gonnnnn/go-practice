package main

import "net/http"

func main(){
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/select", selectHandler)
	http.HandleFunc("/insert", insertHandler)
	http.HandleFunc("/delete", deleteHandler)
	http.ListenAndServe(":8000", nil)

}

func rootHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte("hi"))
}

func selectHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte("selectHandler was called"))
}

func insertHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte("insertHandler was called"))
}

func deleteHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte("deleteHandler was called"))
}