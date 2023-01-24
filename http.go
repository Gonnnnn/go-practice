package main

import "net/http"

func HandleStatic(fileName string) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		http.ServeFile(writer, request, fileName)
	}
}