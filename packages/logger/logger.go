package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Logger struct {
	mu         sync.Mutex
	file       *os.File
	dir        string
	prefix     string
	useConsole bool
	currentDay string
}

var (
	appLogger *Logger
)

func NewLogger(dir, prefix string, useConsole bool) *Logger {
	logger := &Logger{
		dir:        dir,
		prefix:     prefix,
		useConsole: useConsole,
		currentDay: time.Now().Format("02"),
	}

	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		panic(fmt.Sprintf("Failed to create log directory: %v", err))
	}

	logger.createLogFile(false)
	return logger
}

func (l *Logger) createLogFile(isNewDay bool) {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	filename := fmt.Sprintf("%s_%s_%02d%02d%02d.log", l.prefix, now.Format("20060102"), now.Hour(), now.Minute(), now.Second())
	filePath := filepath.Join(l.dir, filename)

	file, err := os.Create(filePath)
	if err != nil {
		fmt.Printf("ERROR: Cannot create log file: %v\n", filePath)
		l.file = nil
		return
	}

	l.file = file

	timestamp := now.Format("20060102 150405 ")
	l.file.WriteString("====================================================================\n")
	if isNewDay {
		l.file.WriteString(timestamp + "BEGIN LOG for a NEW DAY\n")
	} else {
		l.file.WriteString(timestamp + "BEGIN LOG\n")
	}
	l.file.WriteString("====================================================================\n")
}

func (l *Logger) writeLogEntry(entry string) {
	now := time.Now()
	timestamp := now.Format("20060102 150405 .000 ")

	if l.useConsole {
		fmt.Println(timestamp + entry)
	}

	if l.file == nil {
		return
	}

	if now.Format("02") != l.currentDay {
		l.file.WriteString("====================================================================\n")
		l.file.WriteString(timestamp + "END LOG for the OLD DAY\n")
		l.file.WriteString("====================================================================\n")
		l.file.Close()
		l.createLogFile(true)
		l.currentDay = now.Format("02")
	}

	l.file.WriteString(timestamp + entry + "\n")
}

func (l *Logger) WriteLine(v ...interface{}) {
	entry := fmt.Sprint(v...)
	l.writeLogEntry(entry)
}

func (l *Logger) Close() {
	if l.file != nil {
		l.file.Close()
	}
}

func Log(v ...any) {
	if appLogger == nil {
		appLogger = NewLogger("./log", "log", true)
		// defer appLogger.Close()
	}
	appLogger.WriteLine(v...)
}
