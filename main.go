package main

import (
	"chat_server/router"
	"os"

	"github.com/shiena/ansicolor"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.StandardLogger().SetFormatter(&logrus.TextFormatter{ForceColors: true})
	logrus.StandardLogger().SetOutput(ansicolor.NewAnsiColorWriter(os.Stdout))

	r := router.InitRouter()
	panic(r.Run("127.0.0.1:8080"))
}
