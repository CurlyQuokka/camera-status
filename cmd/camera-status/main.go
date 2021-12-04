package main

// ToDo:
// 1. Add logfile
// 2. Parameters/config file
// 3. Delete data when disk when running out of space
// 4. Try to restart daemon when files are not up to date

import (
	"net/http"
	"os"

	"github.com/CurlyQuokka/camera-status/pkg/mailer"
	"github.com/CurlyQuokka/camera-status/pkg/watcher"
	"github.com/CurlyQuokka/camera-status/pkg/webserver"
)

var watcherObj *watcher.Watcher

const (
	mailFrom = "@gmail.com"
	mailTo   = "@gmail.com"
	mailPass = ""

	smtpHost = "smtp.gmail.com"
	smtpPort = "587"
)

func main() {
	var webServerObj *webserver.WebServer

	var port string
	if len(os.Args) < 2 {
		port = ":80"
	} else {
		port = ":" + os.Args[1]
	}

	mailerObj := mailer.NewMailer(mailFrom, mailTo, mailPass, smtpHost, smtpPort)
	watcherObj = watcher.NewDefaultWatcher(mailerObj)
	webServerObj = webserver.NewWebServer(port, watcherObj)

	finished := make(chan bool)

	go watcherObj.Watch(finished)

	http.HandleFunc("/", webServerObj.GetData)

	http.ListenAndServe(port, nil)

	<-finished
}
