package main

import (
	"log"

	"github.com/l-vitaly/goenv/example/config"
)

func main() {
	cfg, err := config.Parse()
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(cfg.DBConn)
}
