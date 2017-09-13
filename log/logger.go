package log

import (
	"fmt"
	"io"
	"os"
	"time"
)

type Logger interface {
	D(string, ...interface{})
	I(string, ...interface{})
	W(string, ...interface{})
	E(string, ...interface{})
}

type LogWriter interface {
	Write(info *LogInfo)
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
	Writer LogWriter
}

type LogInfo struct {
	time  string
	level string
	msg   string
}

var defaultConfig *Config = &Config{
	Level:  DEBUG,
	Cache:  1000,
	Writer: StdWriter,
}

type proxy struct {
	config *Config
	ch     chan *LogInfo
}

func (p *proxy) write() {
	for msg := range p.ch {
		p.config.Writer.Write(msg)
	}
}

func (p *proxy) formatLog(format string, args ...interface{}) {
	now := time.Now()
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
	p.ch <- &logInfo{
		time:  now.Format("2006-01-02 15:04:05"),
		level: level,
		msg: fmt.Sprintf("%s %s: %s",
			now.Format("2006-01-02 15:04:05"),
			level, fmt.Sprintf(format, args...)),
	}

}

func (p *proxy) D(format string, args ...interface{}) {
	if p.config.Level >= DEBUG {
		p.formatLog(format, args...)
	}
}

func (p *proxy) I(format string, args ...interface{}) {
	if p.config.Level >= INFO {
		p.formatLog(format, args...)
	}
}

func (p *proxy) W(format string, args ...interface{}) {
	if p.config.Level >= WARN {
		p.formatLog(format, args...)
	}
}

func (p *proxy) E(format string, args ...interface{}) {
	if p.config.Level >= ERROR {
		p.formatLog(format, args...)
	}
}

func NewLogger(config *Config) Logger {
	if config == nil {
		config = defaultConfig
	}
	logger := &proxy{
		config: config,
		ch:     make(chan *LogInfo, config.Cache),
	}
	go func() {
		logger.write()
	}()
	return logger
}
