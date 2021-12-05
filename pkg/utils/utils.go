package utils

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"os/exec"
	"strings"
	"time"

	"golang.org/x/sys/unix"
)

const (
	base           = 1024
	power          = 2
	spaceThreshold = 0.1

	systemctlCmd = "/usr/bin/systemctl"
)

type File struct {
	Name string
	Size float64
}

type FileList []File

func ListFiles(path string) FileList {
	files, err := ioutil.ReadDir(path)
	result := FileList{}
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}
	for _, file := range files {
		var f File
		f.Name = file.Name()
		f.Size = float64(file.Size()) / math.Pow(base, power)
		result = append(result, f)
	}
	result = result.removeIrrelevant()
	return result
}

func (f File) Print() {
	fmt.Printf("%s - %.2f MB\n", f.Name, f.Size)
}

func (list FileList) removeIrrelevant() FileList {
	var newList FileList
	for _, f := range list {
		if strings.Contains(f.Name, ".mkv") {
			newList = append(newList, f)
		}
	}
	return newList
}

func (list FileList) Revert() FileList {
	var newList FileList
	for i := len(list) - 1; i >= 0; i-- {
		newList = append(newList, list[i])
	}
	return newList
}

func (list FileList) Print() {
	for _, f := range list {
		f.Print()
	}
}

func (list FileList) GetLatest() FileList {
	if len(list) < 4 {
		return list
	}
	return list[:3]
}

func (f File) ParseDate() time.Time {
	now := time.Now().Local()
	timeOffset := "+02:00"
	if !now.IsDST() {
		timeOffset = "+01:00"
	}
	nameSplit := strings.Split(f.Name, ".")
	dateSplit := strings.Split(nameSplit[0], "T")
	dateSplit[1] = strings.ReplaceAll(dateSplit[1], "-", ":")
	dateToParse := dateSplit[0] + "T" + dateSplit[1] + timeOffset

	parsedDate, err := time.Parse(time.RFC3339, dateToParse)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(555)
	}

	return parsedDate
}

func IsUpToDate(date time.Time) bool {
	now := time.Now().Local()
	diff := now.Sub(date)
	return diff.Minutes() <= 1.5
}

func GetFsSpace(path string) float64 {
	var stat unix.Statfs_t
	unix.Statfs(path, &stat)
	return (float64(stat.Bavail) / float64(stat.Blocks))
}

func IsSpaceSufficient(path string) bool {
	return GetFsSpace(path) > spaceThreshold
}

func IsDaemonActive(name string) bool {
	cmd := exec.Command(systemctlCmd, "is-active", "--quiet", name)
	err := cmd.Run()
	return err == nil
}

func RestartDaemon(name string) error {
	cmd := exec.Command(systemctlCmd, "restart", name)
	err := cmd.Run()
	return err
}

func RemoveFiles(path string, numOfFiles int) error {
	files := ListFiles(path)
	for i := 0; i < numOfFiles; i++ {
		if i < len(files) {
			err := os.Remove(path + "/" + files[i].Name)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
