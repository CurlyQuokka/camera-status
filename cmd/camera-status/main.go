package main

// ToDo:
// 1. Add logfile
// 2. Parameters/config file
// 3. Delete data when disk when running out of space
// 4. Try to restart daemon when files are not up to date

import (
	"fmt"
	"net/http"
	"os"

	"github.com/CurlyQuokka/camera-status/pkg/mailer"
	"github.com/CurlyQuokka/camera-status/pkg/utils"
	"github.com/CurlyQuokka/camera-status/pkg/watcher"
)

var watcherObj *watcher.Watcher

const (
	StatusOK  = "OK"
	StatusBad = "BAD"

	green = "#cce36f"
	red   = "#e36f6f"

	mailFrom = "@gmail.com"
	mailTo   = "@gmail.com"
	mailPass = ""

	smtpHost = "smtp.gmail.com"
	smtpPort = "587"
)

func main() {
	mailerObj := mailer.NewMailer(mailFrom, mailTo, mailPass, smtpHost, smtpPort)
	watcherObj = watcher.NewDefaultWatcher(mailerObj)

	var port string

	if len(os.Args) < 2 {
		port = ":80"
	} else {
		port = ":" + os.Args[1]
	}

	finished := make(chan bool)

	go watcherObj.Watch(finished)

	http.HandleFunc("/", getData)

	http.ListenAndServe(port, nil)

	<-finished
}

func getData(w http.ResponseWriter, req *http.Request) {
	latest, isDaemonActive, isUpToDate, isSpaceSufficient, space := watcherObj.CheckStatus()
	writeHTML(isDaemonActive, isUpToDate, isSpaceSufficient, space, latest, w)
}

func writeHTML(isDaemonActive, isUpToDate, isSpaceSufficient bool, space float64, latest utils.FileList, w http.ResponseWriter) {
	bgColor := red
	status := StatusBad

	if isDaemonActive && isUpToDate && isSpaceSufficient {
		bgColor = green
		status = StatusOK
	}

	fmt.Fprintf(w, "<body bgcolor = \"%s\">", bgColor)
	fmt.Fprintf(w, "<h1>STATUS: %s</h1>", status)

	if status == StatusBad {
		if !isDaemonActive {
			fmt.Fprintf(w, "<h2>DAEMON IS INACTIVE!</h2>")
		}
		if !isUpToDate {
			fmt.Fprintf(w, "<h2>RECORDINGS ARE NOT UP TO DATE!</h2>")
		}
		if !isSpaceSufficient {
			fmt.Fprintf(w, "<h2>RUNNING OUT OF SPACE!</h2>")
		}
	}

	for _, f := range latest {
		fmt.Fprintf(w, "%s - %.2f MB</br>", f.Name, f.Size)
	}

	fmt.Fprintf(w, "</br>")
	fmt.Fprintf(w, "Free space: <b>%.2f%%</b></br>", space*100)
	fmt.Fprintf(w, "</body>")
}
