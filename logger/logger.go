package logger

import (
	"fmt"
	"time"
)

type logger struct {
	env        string
	token      string
	serverRoot string
}

var dt time.Time

func init() {
	dt = time.Now()
}

type Logger interface {
	Warn(args ...interface{})
	Info(args ...interface{})
	Error(args ...interface{})
}

func (logger) Warn(args ...interface{}) {
	fmt.Printf("[WARN] %s - %s\n", append([]interface{}{dt.Format("01-02-2006 15:04:05")}, args...)...)
}

func (logger) Info(args ...interface{}) {
	fmt.Printf("[INFO] %s - %s\n", append([]interface{}{dt.Format("01-02-2006 15:04:05")}, args...)...)
}

func (logger) Error(args ...interface{}) {
	fmt.Printf("[Error] %s - %s\n", append([]interface{}{dt.Format("01-02-2006 15:04:05")}, args...)...)
}

func NewLogger(env, token string) Logger {
	return &logger{env: env, token: token, serverRoot: "/"}
}
