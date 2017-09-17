package log

import (
	"io"
)

type StdWriter struct {
	writer io.Writer
}

func (sw *StdWriter) Write(info *LogInfo) {
	sw.writer.Write([]byte(info.msg))
}
