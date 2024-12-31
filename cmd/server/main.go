package main

import (
	"fmt"
	"log"
	"os"
	"net/http"

	"github.com/blobfish465/common-circle-web-forum/internal/router"
)

func main() {
	r := router.Setup()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000" // Default to port 8000 if PORT is not set
		// to ensure backend still works locally when running outside Render.
	}

	fmt.Printf("Listening on port %s at http://0.0.0.0:%s!\n", port, port)
	log.Fatalln(http.ListenAndServe(":"+port, r))
}
