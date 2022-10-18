package logger

import "fmt"

type Logger struct {
	logs []string
}

func NewLogger() *Logger {
	return &Logger{}
}

func (l *Logger) CreateLog(data interface{}) *[]string {
	l.logs = append(l.logs, fmt.Sprintf("%v", data))
	return &l.logs
}

func (l * Logger) GetLogs() []string {
	switch true {
	case l.logs == nil:
		l.logs = append(l.logs, "")
	}
	return l.logs
}