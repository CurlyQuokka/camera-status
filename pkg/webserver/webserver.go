package webserver

import (
	"fmt"
	"net/http"

	"github.com/CurlyQuokka/camera-status/pkg/utils"
	"github.com/CurlyQuokka/camera-status/pkg/watcher"
)

const (
	statusOK  = "OK"
	statusBad = "BAD"

	green = "#cce36f"
	red   = "#e36f6f"
)

type WebServer struct {
	port  string
	watch *watcher.Watcher
}

func NewWebServer(port string, w *watcher.Watcher) *WebServer {
	return &WebServer{
		port:  port,
		watch: w,
	}
}

func (ws *WebServer) GetData(w http.ResponseWriter, req *http.Request) {
	latest, isDaemonActive, isUpToDate, isSpaceSufficient, space := ws.watch.CheckStatus()
	writeHTML(isDaemonActive, isUpToDate, isSpaceSufficient, space, latest, w)
}

func writeHTML(isDaemonActive, isUpToDate, isSpaceSufficient bool, space float64, latest utils.FileList, w http.ResponseWriter) {
	bgColor := red
	status := statusBad

	if isDaemonActive && isUpToDate && isSpaceSufficient {
		bgColor = green
		status = statusOK
	}

	fmt.Fprintf(w, "<body bgcolor = \"%s\">", bgColor)
	fmt.Fprintf(w, "<h1>STATUS: %s</h1>", status)

	if status == statusBad {
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
