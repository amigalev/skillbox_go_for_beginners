package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/goombaio/namegenerator"
)

func main() {
	http.HandleFunc("/", handler)
	fmt.Println(http.ListenAndServe("localhost:8080", nil))
}

func handler(writer http.ResponseWriter, request *http.Request) {
	name := request.URL.Path[1:]
	if name == "" {
		seed := time.Now().UTC().Unix()
		nameG := namegenerator.NewNameGenerator(seed)
		fmt.Fprintf(writer, "Hello, %s", nameG.Generate())
	} else {
		fmt.Fprintf(writer, "Hello, %s", name)
	}
}
