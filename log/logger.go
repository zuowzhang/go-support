package log

import (
	"os"
	"io"
	"fmt"
	"time"
)

type Logger interface {
	D(string, ...interface{})
	I(string, ...interface{})
	W(string, ...interface{})
	E(string, ...interface{})
}

const (
	NULL int = iota
	DEBUG
	INFO
	WARN
	ERROR
)

type Config struct {
	Level  int
	Cache  int
	Writer io.Writer
}

var defaultConfig *Config = &Config{
	Level:DEBUG,
	Cache:200,
	Writer:os.Stdout,
}

type proxy struct {
	config *Config
	ch     chan string
}

func (p *proxy)write() {
	for msg := range p.ch {
		p.config.Writer.Write([]byte(msg))
	}
}

func (p *proxy)formatLog(format string, args ...interface{}) {
	var level string
	switch p.config.Level {
	case DEBUG:
		level = "DEBUG"
	case INFO:
		level = "INFO"
	case WARN:
		level = "WARN"
	default:
		level = "ERROR"

	}
	p.ch <- fmt.Sprintf("%s %s: %s",
		time.Now().Format("2006-01-02 15:04:05"),
		level, fmt.Sprintf(format, args...))

}

func (p *proxy)D(format string, args ...interface{}) {
	if p.config.Level >= DEBUG {
		p.formatLog(format, args...)
	}
}

func (p *proxy)I(format string, args ...interface{}) {
	if p.config.Level >= INFO {
		p.formatLog(format, args...)
	}
}

func (p *proxy)W(format string, args ...interface{}) {
	if p.config.Level >= WARN {
		p.formatLog(format, args...)
	}
}

func (p *proxy)E(format string, args ...interface{}) {
	if p.config.Level >= ERROR {
		p.formatLog(format, args...)
	}
}

func NewLogger(config *Config) Logger {
	if config == nil {
		config = defaultConfig
	}
	logger := &proxy{
		config:config,
		ch:make(chan string, config.Cache),
	}
	go func() {
		logger.write()
	}()
	return logger
}
