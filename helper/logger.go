package helper

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"time"
)

func LoggerError(keyword, data string) {
	if os.Getenv("APP_MODE") == "production" {
		if _, err := os.Stat("logs/error"); os.IsNotExist(err) {
			fmt.Println(os.MkdirAll("logs/error", 0700))
		}
		newData := fmt.Sprintf("%s %s", keyword, data)
		loggerLog("logs/error/", time.Now().Format("2006-01-02"), newData)
	} else {
		fmt.Println(data)
	}
}

func LoggerAccees(data string) {
	if os.Getenv("APP_MODE") == "production" {
		if _, err := os.Stat("logs/access"); os.IsNotExist(err) {
			fmt.Println(os.MkdirAll("logs/access", 0700))
		}
		loggerLog("logs/access/", time.Now().Format("2006-01-02"), data)
	} else {
		fmt.Println(data)
	}
}

func loggerLog(dirname, filename, data string) {
	d := createFile(dirname, filename)
	f, err := os.OpenFile(dirname+d.Name(), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf(err.Error())
	}
	defer f.Close()
	log.SetOutput(f)
	log.SetFlags(0)
	tm := fmt.Sprintf("[%s] ", time.Now().Format("2006-01-02 15:04:05"))
	log.SetPrefix(tm)
	log.Println(data)
}

func createFile(dirname, name string) fs.FileInfo {
	f, err := os.Stat(fmt.Sprintf("%s%s.log", dirname, name))
	if err != nil {
		os.Create(fmt.Sprintf("%s%s.log", dirname, name))
		return createFile(dirname, name)
	}
	return f
}
