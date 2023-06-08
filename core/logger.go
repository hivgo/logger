package main

import (
	"fmt"
	ct "github.com/daviddengcn/go-colortext"
	"os"
	"sync"
	"time"
)

const (
	MaxLogLine = 65535
	LogLevel   = LTInfo
)

const (
	LTInvalid = iota
	LTDebug
	LTInfo
	LTWarn
	LTError
)

type LogFile struct {
	file  *os.File
	lines int
	path  string
	mutex sync.Mutex
	time  time.Time
	num   int
}

// Log 是一个写日志文件的函数
func (log *LogFile) Log(text string) {
	log.mutex.Lock()
	defer log.mutex.Unlock()

	time := time.Now()
	var dir string

	if log.time.Day() != time.Day() {
		log.Close()
	}

	if log.lines >= MaxLogLine {
		if createFile(dir, time, log) {
			return
		}
		log.lines = 0
	}

	if log.file == nil {
		dir = fmt.Sprintf(log.path+"/%.4d_%.2d%.2d", time.Year(), time.Month(), time.Day())
		os.MkdirAll(dir, os.ModeDir)

		if createFile(dir, time, log) {
			return
		}
	}
	log.file.WriteString(text + " \n")
	//
	//io.WriteString(log.file,text)
	log.lines++
}

func createFile(dir string, time time.Time, log *LogFile) bool {
	filePath := fmt.Sprintf(dir+"/%.4d_%.2d%.2d_%.1d.log", time.Year(), time.Month(), time.Day(), log.num)
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Printf("Create log file %s error faild:%s", filePath, err)
		return true
	}

	log.file = file
	log.time = time
	log.num++
	return false
}

// Close 还原LogFile
func (log *LogFile) Close() {
	log.file.Close()
	log.file = nil
	log.lines = 0
	log.num = 1
}

var _sysLog = &LogFile{path: "logs"}
var sysMutex sync.Mutex

func writeConsole(color ct.Color, text string) {
	sysMutex.Lock()
	defer sysMutex.Unlock()

	ct.ChangeColor(color, true, ct.Black, false)
	fmt.Println(text)
	ct.ResetColor()
}

func write(console bool, logType int, format string, args ...interface{}) {
	if logType < LogLevel {
		return
	}

	logTypeKey := "U"
	color := ct.White
	switch logType {
	case LTDebug:
		logTypeKey = "DEBUG"
		color = ct.White
	case LTInfo:
		logTypeKey = "INFO"
		color = ct.Green
	case LTError:
		logTypeKey = "ERROR"
		color = ct.Red
	case LTWarn:
		logTypeKey = "WARN"
		color = ct.Yellow
	}

	time := time.Now()
	var logFormat = fmt.Sprintf("[%s] [%s] %s",
		time.Format("2006-01-02 15:04:05"), logTypeKey, format)
	var text = fmt.Sprintf(logFormat, args...)
	if console {
		writeConsole(color, text)
	}
	_sysLog.Log(text)

}

func Error(format string, args ...interface{}) {
	write(true, LTError, format, args)
}

func Info(format string, args ...interface{}) {
	write(true, LTInfo, format, args)
}

func Warn(format string, args ...interface{}) {
	write(true, LTWarn, format, args)
}

func Debug(format string, args ...interface{}) {
	write(true, LTDebug, format, args)
}

func Log(format string, args ...interface{}) {
	write(false, LTInfo, format, args)
}

func main() {
	Error("Test : %s", "receive handle query list guild params roleId:10003")
}
