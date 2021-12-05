package main

// ToDo:
// 1. Parameters/config file

import (
	"net/http"
	"os"

	"github.com/CurlyQuokka/camera-status/pkg/logger"
	"github.com/CurlyQuokka/camera-status/pkg/mailer"
	"github.com/CurlyQuokka/camera-status/pkg/watcher"
	"github.com/CurlyQuokka/camera-status/pkg/webserver"
)

const (
	mailFrom = "@gmail.com"
	mailTo   = "@gmail.com"
	mailPass = ""

	smtpHost = "smtp.gmail.com"
	smtpPort = "587"
)

func main() {
	loggerObj := logger.NewMyLogger()

	var webServerObj *webserver.WebServer

	var port string
	if len(os.Args) < 2 {
		port = ":80"
	} else {
		port = ":" + os.Args[1]
	}

	mailerObj := mailer.NewMailer(mailFrom, mailTo, mailPass, smtpHost, smtpPort, loggerObj)
	watcherObj := watcher.NewDefaultWatcher(mailerObj, loggerObj)
	webServerObj = webserver.NewWebServer(port, watcherObj, loggerObj)

	finished := make(chan bool)

	go watcherObj.Watch(finished)

	loggerObj.Info("HTTP webserver will start at port " + port)
	http.HandleFunc("/", webServerObj.GetData)
	http.ListenAndServe(port, nil)

	<-finished
}
