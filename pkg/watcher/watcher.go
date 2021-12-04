package watcher

import (
	"fmt"
	"time"

	"github.com/CurlyQuokka/camera-status/pkg/mailer"
	"github.com/CurlyQuokka/camera-status/pkg/utils"
)

const (
	defaultRecDaemon = "camera-recording"
	defaultRecDir    = "/home/user/camera-disk"

	defaultCheckInterval = 10
	numOfFilesPerDay     = 1440
)

type Watcher struct {
	StatusCheckInterval int
	StatusSlice         []bool
	Path                string
	RecordingDaemon     string
	mailer              *mailer.Mailer
}

func NewDefaultWatcher(m *mailer.Mailer) *Watcher {
	return NewCustomWatcher(defaultCheckInterval, defaultRecDir, defaultRecDaemon, m)
}

func NewCustomWatcher(statusInterval int, directory, recordingDaemon string, mail *mailer.Mailer) *Watcher {
	w := &Watcher{
		StatusCheckInterval: statusInterval,
		StatusSlice:         []bool{},
		Path:                directory,
		RecordingDaemon:     recordingDaemon,
		mailer:              mail,
	}
	return w
}

func (w *Watcher) Watch(finished chan bool) {
	for {
		lastStatusLen := len(w.StatusSlice)
		_, isDaemonActive, isUpToDate, isSpaceSufficient, space := w.CheckStatusAndUpdate()

		if !isDaemonActive || !isUpToDate {
			fmt.Println("Daemon will be restarted")
			err := utils.RestartDaemon(w.RecordingDaemon)
			if err != nil {
				fmt.Println(err.Error())
			} else {
				fmt.Println("Daemon successfully restarted")
			}
			time.Sleep(time.Duration(w.StatusCheckInterval) * time.Second)
			_, isDaemonActive, isUpToDate, isSpaceSufficient, space = w.CheckStatusAndUpdate()
		}

		if !isSpaceSufficient {
			fmt.Println("Will remove files now")
			err := utils.RemoveFiles(w.Path, numOfFilesPerDay/2)
			if err != nil {
				fmt.Println(err.Error())
			} else {
				fmt.Printf("Removed %d files", numOfFilesPerDay/2)
				_, isDaemonActive, isUpToDate, isSpaceSufficient, space = w.CheckStatusAndUpdate()
			}
		}

		currentStatusLen := len(w.StatusSlice)

		errorMsg := ""

		if lastStatusLen != currentStatusLen {
			if currentStatusLen > 0 {
				if !isDaemonActive {
					errorMsg += daemonInactiveMsg()
				}
				if !isUpToDate {
					errorMsg += notUpToDateMsg()
				}
				if !isSpaceSufficient {
					errorMsg += spaceMsg(space)
				}
				err := w.mailer.SendMail(mailer.ErrSubject, errorMsg)
				if err != nil {
					fmt.Println(err.Error())
				}
				fmt.Println(errorMsg)
				fmt.Printf("Send error mail. lastStatusLen: %d currentStatusLen: %d\n", lastStatusLen, currentStatusLen)
			} else {
				err := w.mailer.SendMail(mailer.OkSubject, "Everything is OK!\r\n")
				if err != nil {
					fmt.Println(err.Error())
				}
				fmt.Printf("Send ok mail. lastStatusLen: %d currentStatusLen: %d\n", lastStatusLen, currentStatusLen)
			}
		}

		time.Sleep(time.Duration(w.StatusCheckInterval) * time.Second)
	}
}

func (w *Watcher) CheckStatusAndUpdate() (utils.FileList, bool, bool, bool, float64) {
	w.StatusSlice = []bool{}

	latest, isDaemonActive, isUpToDate, isSpaceSufficient, space := w.CheckStatus()

	w.processStatus(isDaemonActive)
	w.processStatus(isUpToDate)
	w.processStatus(isSpaceSufficient)

	return latest, isDaemonActive, isUpToDate, isSpaceSufficient, space

}

func (w *Watcher) CheckStatus() (utils.FileList, bool, bool, bool, float64) {
	files := utils.ListFiles(w.Path)
	files = files.Revert()
	latest := files.GetLatest()
	date := latest[0].ParseDate()

	isDaemonActive := utils.IsDaemonActive(w.RecordingDaemon)
	isUpToDate := utils.IsUpToDate(date)
	space := utils.GetFsSpace(w.Path)
	isSpaceSufficient := utils.IsSpaceSufficient(w.Path)

	return latest, isDaemonActive, isUpToDate, isSpaceSufficient, space

}

func spaceMsg(space float64) string {
	return "The space is running out! Available: " + fmt.Sprintf("%.2f", space*100) + "\r\n"
}

func notUpToDateMsg() string {
	return "Recordings are not up-to-date!\r\n"
}

func daemonInactiveMsg() string {
	return "Camera daemon is inactive!\r\n"
}

func (w *Watcher) processStatus(status bool) {
	if !status {
		w.StatusSlice = append(w.StatusSlice, status)
	}
}
