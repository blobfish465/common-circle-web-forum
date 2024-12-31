package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/blobfish465/common-circle-web-forum/internal/router"
)

func main() {
	r := router.Setup()

	fmt.Print("Listening on port 8080 at http://0.0.0.0:8080!")

	log.Fatalln(http.ListenAndServe("0.0.0.0:8080", r))
}
