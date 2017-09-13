package log

import (
	"io"
	"os"
)

type StdWriter struct {
	writer io.Writer
}

func (sw *StdWriter) Write(info *LogInfo) {

}
