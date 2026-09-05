package main

import (
	"log"

	"github.com/jevvonn/ws-cicd-bcc2026/internal/bootstrap"
)

func main() {
	if err := bootstrap.Run(); err != nil {
		log.Fatal(err)
	}
}
