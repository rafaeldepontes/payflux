package main

import (
	"log"
	"net/http"

	"github.com/rafaeldepontes/goplo/internal/handler"
)

func main() {
	h := handler.NewHandler()

	log.Fatalln(http.ListenAndServe(":8080", h))
}
