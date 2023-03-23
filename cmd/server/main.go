package main

import (
	"fmt"
	"log"

	"github.com/Ragnar-BY/companies-handler/internal/config"
)

func main() {
	cfg, err := config.LoadConfig(".env")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(cfg)
}
