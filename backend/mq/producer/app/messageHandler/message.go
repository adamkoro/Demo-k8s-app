package messageHandler

import (
	"log"
	"os"
)

var (
	WarningLogger = log.New(os.Stdout, "[WARNING] ", log.Ldate|log.Ltime)
	ErrorLogger   = log.New(os.Stderr, "[ERROR] ", log.Ldate|log.Ltime)
	InfoLogger    = log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime)
)

func IsError(err error) bool {
	return err != nil
}
