package log

import (
	"os"
	"sync"
)

type message interface {
	String([]byte) []byte
}

type MessageLog struct {
	mux  sync.Mutex
	name string
	file *os.File
	q    chan message
	wg   sync.WaitGroup
}

func (l *MessageLog) Open(name string) error {
	file, err := os.OpenFile(name, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	l.name = name
	l.file = file
	l.q = make(chan message, 10000)
	l.wg.Add(1)
	go l.save()
	return nil
}

func (l *MessageLog) Close() {
	close(l.q)
	l.wg.Wait()
	l.file.Close()
}

func (l *MessageLog) Write(msg message) {
	l.q <- msg
}

func (l *MessageLog) Switch() error {
	file, err := os.OpenFile(l.name, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	l.mux.Lock()
	defer l.mux.Unlock()
	l.file.Close()
	l.file = file
	return nil
}

func (l *MessageLog) save() {
	defer l.wg.Done()
	buffer := make([]byte, 0, 32*1024)
	for m := range l.q {
		buf := buffer[:0]
		buf = m.String(buf)
		buf = append(buf, '\n')
		l.mux.Lock()
		_, _ = l.file.Write(buf)
		l.mux.Unlock()
	}
}
