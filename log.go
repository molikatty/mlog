package mlog

import (
	"errors"
	"os"
	"strings"
	"sync"
)

type writer interface {
	WriteString(s string) (i int, err error)
}

type Log struct {
	Stdout writer
	Stderr writer
}

var (
	ErrOpenFile = errors.New("open file error")
)

var (
	logger *Log
	once   sync.Once
)

func Logger() *Log {
	if logger == nil {
		once.Do(func() {
			logger = &Log{
				Stdout: os.Stdout,
				Stderr: os.Stderr,
			}
		})
	}

	return logger
}

func (log *Log) Err(perfix, msg string) {
	log.OutMsg(log.Stderr, perfix, msg)
}

func (log *Log) Info(perfix, msg string) {
	log.OutMsg(log.Stdout, perfix, msg)
}

func (log *Log) OutMsg(w writer, perfix, msg string) {
	var s strings.Builder
	s.WriteString(perfix)
	s.WriteString(" ")
	s.WriteString(msg)
	w.WriteString(s.String())
}

func (log *Log) WriteFile(name string, msg ...string) error {
	file, err := os.OpenFile(name, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	for i := range msg {
		log.OutMsg(file, "", msg[i])
	}

	return nil
}
