package log

import (
	"os"
	"io/ioutil"
	"strings"
	"sort"
	"time"
)

const log_suffix string = ".log"

type fileWriter struct {
	dir             string
	count           int
	dirExists       bool
	countValid      bool
	pFile           *os.File
	currentFileName string
}

func NewFileWriter(dir string, fileCount int) LogWriter {
	if dir == "" {
		dir = "."
	}
	if fileCount < 1 {
		fileCount = 7
	}
	return &fileWriter{
		dir:dir,
		count:fileCount,
	}
}

func (fw *fileWriter)checkDir() bool {
	if !fw.dirExists {
		if fw.dir == "" {
			fw.dir = "."
			fw.dirExists = true
		} else {
			_, err := os.Stat(fw.dir)
			if err == nil {
				fw.dirExists = true
			} else if os.IsNotExist(err) {
				if os.MkdirAll(fw.dir, 0777) == nil {
					fw.dirExists = true
				}
			}
		}
	}
	return fw.dirExists
}

type fileNameAndModifyTime struct {
	name       string
	lastModify time.Time
}

type fileNameAndModifyTimeSlice []fileNameAndModifyTime

func (s fileNameAndModifyTimeSlice)Len() int {
	return len(s)
}

func (s fileNameAndModifyTimeSlice)Less(i, j int) bool {
	return s[i].lastModify.Unix() < s[j].lastModify.Unix()
}

func (s fileNameAndModifyTimeSlice)Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (fw *fileWriter)checkFileCount() {
	fileInfos, err := ioutil.ReadDir(fw.dir)
	if err != nil {
		return
	}
	var logFileInfos fileNameAndModifyTimeSlice
	for _, fileInfo := range fileInfos {
		if !fileInfo.IsDir() && strings.HasSuffix(fileInfo.Name(), log_suffix) {
			logFileInfos = append(logFileInfos, fileNameAndModifyTime{
				name:fileInfo.Name(),
				lastModify:fileInfo.ModTime(),
			})
		}
	}
	if len(logFileInfos) > fw.count {
		//先对日志文件进行排序
		sort.Sort(logFileInfos)
		//删除超过最大文件数的日志文件
		for i := 0; i < fw.count; i++ {
			os.Remove(fw.dir + string(os.PathSeparator) + logFileInfos[i].name)
		}
	}
}

func (fw *fileWriter)openLogFile(fileName string) error {
	pFile, err := os.OpenFile(fw.dir + string(os.PathSeparator) + fileName + log_suffix,
		os.O_CREATE | os.O_APPEND | os.O_RDWR,
		0666)
	if err != nil {
		return err
	}
	fw.pFile = pFile
	fw.currentFileName = fileName
	fw.checkFileCount()
	return nil
}

func (fw *fileWriter) checkFileName(fileName string) error {
	if fw.pFile == nil {
		err := fw.openLogFile(fileName)
		if err != nil {
			return err
		}
	} else {
		if fw.currentFileName != fileName {
			fw.pFile.Close()
			fw.pFile = nil
			err := fw.openLogFile(fileName)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (fw *fileWriter) Write(info *LogInfo) {
	if fw.checkDir() {
		err := fw.checkFileName(info.time)
		if err != nil {
			return
		}
		fw.pFile.Write([]byte(info.msg))
	}
}
