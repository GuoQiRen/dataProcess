package log

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

const (
	levelDebug = 1
	levelInfo  = 2
	levelWarn  = 3
	levelError = 4
	levelFatal = 5
)

const (
	FTime = 1 << iota
	FLevel
	FLine
)

var logger = NewLogger(FTime | FLevel | FLine)

func NewLogger(flag int) *Logger {
	return &Logger{flag: flag, out: os.Stdout, buf: make([]byte, 32*1024)}
}

func SetLogFile(name string) error          { return logger.SetLogFile(name) }
func Debug(format string, v ...interface{}) { logger.output(levelDebug, fmt.Sprintf(format, v...)) }
func Info(format string, v ...interface{})  { logger.output(levelInfo, fmt.Sprintf(format, v...)) }
func Warn(format string, v ...interface{})  { logger.output(levelWarn, fmt.Sprintf(format, v...)) }
func Error(format string, v ...interface{}) { logger.output(levelError, fmt.Sprintf(format, v...)) }
func Fatal(format string, v ...interface{}) {
	logger.output(levelFatal, fmt.Sprintf(format, v...))
	logger.mux.Lock()
	defer logger.mux.Unlock()
	if logger.file != nil {
		_ = logger.file.Close()
	}
	os.Exit(1)
}

type Logger struct {
	mux  sync.Mutex
	file *os.File
	flag int
	out  io.Writer
	buf  []byte
}

func (l *Logger) SetLogFile(name string) error {
	f, err := os.OpenFile(name, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	l.mux.Lock()
	defer l.mux.Unlock()
	if l.file != nil {
		l.file.Close()
	}
	l.file = f
	l.out = f
	return nil
}

func (l *Logger) Debug(format string, v ...interface{}) {
	l.output(levelDebug, fmt.Sprintf(format, v...))
}

func (l *Logger) Info(format string, v ...interface{}) {
	l.output(levelInfo, fmt.Sprintf(format, v...))
}

func (l *Logger) Warn(format string, v ...interface{}) {
	l.output(levelWarn, fmt.Sprintf(format, v...))
}

func (l *Logger) Error(format string, v ...interface{}) {
	l.output(levelError, fmt.Sprintf(format, v...))
}

func (l *Logger) Fatal(format string, v ...interface{}) {
	l.output(levelFatal, fmt.Sprintf(format, v...))
	l.mux.Lock()
	defer l.mux.Unlock()
	if l.file != nil {
		_ = l.file.Close()
	}
	os.Exit(1)
}

func (l *Logger) output(level int, s string) {
	l.mux.Lock()
	defer l.mux.Unlock()
	buf := l.buf[:0]
	if l.flag > 0 {
		buf = l.formatHeader(buf, level)
	}
	buf = append(buf, s...)
	buf = append(buf, '\n')
	_, _ = l.out.Write(buf)
}

func (l *Logger) formatHeader(buf []byte, level int) []byte {
	space := false
	buf = append(buf, '[')
	if l.flag&FTime != 0 {
		space = true
		buf = AppendTime(buf, time.Now())
	}
	if l.flag&FLevel != 0 {
		if space {
			buf = append(buf, ' ')
		}
		space = true
		switch level {
		case levelDebug:
			buf = append(buf, "DEBUG  "...)
		case levelInfo:
			buf = append(buf, "INFO   "...)
		case levelWarn:
			buf = append(buf, "WARNING"...)
		case levelError:
			buf = append(buf, "ERROR  "...)
		case levelFatal:
			buf = append(buf, "FATAL  "...)
		}
	}
	if l.flag&FLine != 0 {
		if space {
			buf = append(buf, ' ')
		}
		space = true
		_, file, line, ok := runtime.Caller(3)
		if !ok {
			file = "???"
			line = 0
		}
		idx := strings.LastIndex(file, "src/")
		buf = append(buf, file[idx+4:]...)
		buf = append(buf, ':')
		buf = AppendInt(buf, line, -1)
	}
	return append(buf, "] "...)
}

func AppendTime(buf []byte, t time.Time) []byte {
	year, month, day := t.Date()
	buf = AppendInt(buf, year, 2)
	buf = append(buf, '-')
	buf = AppendInt(buf, int(month), 2)
	buf = append(buf, '-')
	buf = AppendInt(buf, day, 2)
	buf = append(buf, 'T')
	hour, min, sec := t.Clock()
	buf = AppendInt(buf, hour, 2)
	buf = append(buf, ':')
	buf = AppendInt(buf, min, 2)
	buf = append(buf, ':')
	buf = AppendInt(buf, sec, 2)
	buf = append(buf, '.')
	buf = AppendInt(buf, t.Nanosecond()/1e6, 3)
	_, offset := t.Zone()
	h := offset / 3600
	if h >= 0 {
		buf = append(buf, '+')
	} else {
		buf = append(buf, '-')
	}
	buf = append(buf, byte('0'+h/10))
	buf = append(buf, byte('0'+h%10))
	m := offset % 3600
	buf = append(buf, byte('0'+m/10))
	return append(buf, byte('0'+m%10))
}

func AppendInt(buf []byte, i int, wid int) []byte {
	var buffer [20]byte
	bp := len(buffer) - 1
	for i >= 10 || wid > 1 {
		q := i / 10
		buffer[bp] = byte('0' + i%10)
		bp--
		wid--
		i = q
	}
	buffer[bp] = byte('0' + i)
	return append(buf, buffer[bp:]...)
}
