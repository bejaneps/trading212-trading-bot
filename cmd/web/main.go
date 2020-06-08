package main

import (
	"os"

	"github.com/bejaneps/trading212/cmd/web/sub"
	log "github.com/sirupsen/logrus"
)

func main() {
	//log.SetReportCaller(true)
	log.SetOutput(os.Stdout)

	if err := sub.Execute(); err != nil {
		log.Fatalln(err)
	}
}
