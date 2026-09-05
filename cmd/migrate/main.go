package main

import (
	"flag"
	"log"

	"github.com/jevvonn/ws-cicd-bcc2026/internal/bootstrap"
	"github.com/jevvonn/ws-cicd-bcc2026/internal/infra/postgresql"
)

func main() {
	mode := flag.String("m", "up", "migration mode: up or down")
	flag.Parse()

	bootstrap.LoadEnv()

	db, err := postgresql.New()
	if err != nil {
		log.Fatal(err)
	}

	switch *mode {
	case "up":
		if err := postgresql.Migrate(db); err != nil {
			log.Fatal(err)
		}
		log.Println("migration up completed")
	case "down":
		if err := postgresql.Rollback(db); err != nil {
			log.Fatal(err)
		}
		log.Println("migration down completed")
	default:
		log.Fatalf("unknown migration mode: %s", *mode)
	}
}
