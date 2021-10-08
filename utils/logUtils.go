package utils

import (
	"fmt"
	"log"
	"os"
)

var FileName string = "audioBookDebug.log"

var (
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
	Debug   *log.Logger
)

func SetupLog() *log.Logger {

	f, err := os.OpenFile(FileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("open file error !")
	}

	fmt.Println("Successfully initialize logger!")
	logger := log.New(f, "[Debug]", log.Ldate|log.Ltime)

	return logger
}
