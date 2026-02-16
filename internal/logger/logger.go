package logger

import (
	"log"
	"os"
)

type Logger struct {
	l *log.Logger
}

func NewApi() *Logger {
	return &Logger{
		l: log.New(os.Stdout, "|API|", log.LstdFlags),
	}
}
func NewDb() *Logger {
	return &Logger{
		l: log.New(os.Stdout, "|DB|", log.LstdFlags),
	}
}
func NewServer() *Logger {
	return &Logger{
		l: log.New(os.Stdout, "|SERVER|", log.LstdFlags),
	}
}
func (l *Logger) PrintErr(err error) {
	l.l.SetFlags(log.LstdFlags | log.Lshortfile)
	l.l.Printf("ERROR: %s", err.Error())
	l.l.SetFlags(log.LstdFlags)
}
func (l *Logger) Print(format string, v ...any) {
	l.l.Printf(format, v...)
}
func (l *Logger) Fatal(format string, v ...any) {
	l.l.Fatalf(format, v...)
}
