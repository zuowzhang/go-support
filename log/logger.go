package log

import (
	"fmt"
	"time"
	"os"
	"sync"
)

type Logger interface {
	D(string, ...interface{})
	I(string, ...interface{})
	W(string, ...interface{})
	E(string, ...interface{})
	CloseSafely()
}

type LogWriter interface {
	Write(info *LogInfo)
}

const (
	DEBUG int = iota
	INFO
	WARN
	ERROR
	NULL
)

type Config struct {
	Level  int
	Cache  int
	Writer LogWriter
}

type LogInfo struct {
	time string
	msg  string
}

var defaultConfig *Config = &Config{
	Level:  DEBUG,
	Cache:  1000,
	Writer: &StdWriter{
		writer:os.Stdout,
	},
}

type proxy struct {
	config    *Config
	ch        chan *LogInfo
	waitGroup sync.WaitGroup
}

func (p *proxy) write() {
	for info := range p.ch {
		p.config.Writer.Write(info)
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
	p.ch <- &LogInfo{
		time:  now.Format("2006-01-02"),
		msg: fmt.Sprintf("%s %s: %s",
			now.Format("2006-01-02 15:04:05"),
			level, fmt.Sprintf(format, args...)),
	}

}

func (p *proxy) D(format string, args ...interface{}) {
	if p.config.Level <= DEBUG {
		p.formatLog(format, args...)
	}
}

func (p *proxy) I(format string, args ...interface{}) {
	if p.config.Level <= INFO {
		p.formatLog(format, args...)
	}
}

func (p *proxy) W(format string, args ...interface{}) {
	if p.config.Level <= WARN {
		p.formatLog(format, args...)
	}
}

func (p *proxy) E(format string, args ...interface{}) {
	if p.config.Level <= ERROR {
		p.formatLog(format, args...)
	}
}

func (p *proxy)CloseSafely() {
	close(p.ch)
	p.waitGroup.Wait()
}

func SetDefaultLevel(level int) {
	defaultConfig.Level = level
}

func setDefaultCache(cache int) {
	defaultConfig.Cache = cache
}

func setDefaultWriter(writer LogWriter) {
	defaultConfig.Writer = writer
}

func NewLogger(config *Config) Logger {
	if config == nil {
		config = defaultConfig
	}
	logger := &proxy{
		config: config,
		ch:     make(chan *LogInfo, config.Cache),
	}
	logger.waitGroup.Add(1)
	go func() {
		logger.write()
		logger.waitGroup.Done()
	}()
	return logger
}
