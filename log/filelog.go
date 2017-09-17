package log

import (
	"os"
	"io/ioutil"
	"strings"
	"sort"
	"time"
)

const log_suffix string = ".log"

type FileWriter struct {
	FileDir         string
	FileCount       int
	dirExists       bool
	pFile           *os.File
	currentFileName string
}

func (fw *FileWriter)checkDir() bool {
	if !fw.dirExists {
		if fw.FileDir == "" {
			fw.FileDir = "."
			fw.dirExists = true
		} else {
			_, err := os.Stat(fw.FileDir)
			if err == nil {
				fw.dirExists = true
			} else if os.IsNotExist(err) {
				if os.MkdirAll(fw.FileDir, 0777) == nil {
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

func (fw *FileWriter)checkFileCount() {
	fileInfos, err := ioutil.ReadDir(fw.FileDir)
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
		} else {
			continue
		}
	}
	if len(logFileInfos) > fw.FileCount {
		//先对日志文件进行排序
		sort.Sort(logFileInfos)
		//删除超过最大文件数的日志文件
		for i := 0; i < fw.FileCount; i++ {
			os.Remove(fw.FileDir + string(os.PathSeparator) + logFileInfos[i].name)
		}
	}
}

func (fw *FileWriter)openLogFile(fileName string) error {
	pFile, err := os.OpenFile(fw.FileDir + string(os.PathSeparator) + fileName + log_suffix,
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

func (fw *FileWriter) checkFileName(fileName string) error {
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

func (fw *FileWriter) Write(info *LogInfo) {
	if fw.checkDir() {
		err := fw.checkFileName(info.time)
		if err != nil {
			return
		}
		fw.pFile.Write([]byte(info.msg))
	}
}
