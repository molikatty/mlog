package mlog

import (
	"errors"
	"os"
	"strings"
	"sync"

	"github.com/wzshiming/ctc"
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
	log.OutMsg(log.Stdout, ctc.ForegroundRed.String()+perfix+ctc.Reset.String(), msg)
}

func (log *Log) Info(perfix, msg string) {
	log.OutMsg(log.Stdout, ctc.ForegroundGreen.String()+perfix+ctc.Reset.String(), msg)
}

func (log *Log) Warning(perfix, msg string) {
	log.OutMsg(log.Stdout, ctc.ForegroundYellow.String()+perfix+ctc.Reset.String(), msg)
}

func (log *Log) OutMsg(w writer, perfix, msg string) {
	var s strings.Builder
	s.Grow(len(perfix) + len(msg) + 1)
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
