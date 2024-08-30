package server

import (
	"log"
	"net/http"
)

func RunWebservice(s *http.Server, err chan error) {

	log.Println("Running HTTP Server")

	err <- s.ListenAndServe()
}
