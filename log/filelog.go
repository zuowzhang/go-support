package log

import (
	"os"
	"strings"
	"time"
)

type FileWriter struct {
	FileDir   string
	FileCount int
	pFile     *os.File
}

func (fw *FileWriter) checkFile() {
	fileName := time.Now().Format("2006-01-02")
	if fw.pFile == nil {
		pFile, err := os.OpenFile(fw.FileDir+os.PathSeparator+fileName,
			os.O_CREATE|os.O_APPEND|os.O_RDWR,
			0666)
		if err != nil {
			return nil
		}
		fw.pFile = pFile
	} else {
		lastIndexOfSeparator := strings.LastIndex(fw.pFile.Name(), os.PathListSeparator)

		fw.pFile.Name()
	}
}

func (fw *FileWriter) Write(p []byte) (n int, err error) {

	return 0, nil
}
