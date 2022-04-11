package main

import (
	"chat_server/router"
	"os"

	"github.com/shiena/ansicolor"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.StandardLogger().SetFormatter(&log.TextFormatter{ForceColors: true})
	log.StandardLogger().SetOutput(ansicolor.NewAnsiColorWriter(os.Stdout))

	r := router.InitRouter()
	panic(r.Run("127.0.0.1:8080"))
}
