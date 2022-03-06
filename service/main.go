package main

import (
	"fmt"
	"os"
	"rplss/service/app"
)

func main() {
	p, _ := os.Getwd()
	app.Run(fmt.Sprintf("%s/service/config/config.yml", p))
}
