package filter

import (
	"github.com/zuowzhang/go-support/simplehttp"
	"github.com/zuowzhang/go-support/log"
	"time"
)

var logger log.Logger

func NewLogFilter(config *log.Config) simplehttp.FilterFunc {
	logger = log.NewLogger(config)
	return func(h simplehttp.HandlerFunc) simplehttp.HandlerFunc {
		return func(c simplehttp.Context) error {
			start := time.Now()
			err := h(c)
			end := time.Now()
			logger.D("%s %s costs %d Milliseconds\n", c.Request().Method, c.Request().URL.Path, end.Sub(start) / 1000)
			return err
		}
	}
}
