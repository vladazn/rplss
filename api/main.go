package main

import (
	"fmt"
	"os"
	"rplss/api/app"
)

func main() {
	p, _ := os.Getwd()
	app.Run(fmt.Sprintf("%s/api/config/config.yml", p))
}
