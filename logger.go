package pig

import (
	"fmt"
	"log"
)

type Logger struct {
}

func NewLogger() *Logger {
	return &Logger{}
}

func (l *Logger) Info(message string) {
	//TODO implement me
	log.Println(fmt.Sprintf("[info] %s", message))
}

func (l *Logger) Debug(message string) {
	//TODO implement me
	log.Println(fmt.Sprintf("[debug] %s", message))
}

func (l *Logger) Warning(message string) {
	//TODO implement me
	log.Println(fmt.Sprintf("[warning] %s", message))
}

func (l *Logger) Fatal(message string) {
	//TODO implement me
	log.Println(fmt.Sprintf("[fatal] %s", message))
}
