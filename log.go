package logrus_log

import (
	"bufio"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/coreservice-io/log"
	"github.com/coreservice-io/logrus_log/nested"
	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
)

var logsAllAbsFolder string
var logsErrorAbsFolder string

type Fields = logrus.Fields
type LogLevel = log.LogLevel

type LocalLog struct {
	*logrus.Logger
	ALL_LogfolderABS string
	ERR_LogfolderABS string
	MaxSize          int
	MaxBackups       int
	MaxAge           int
}

func (logger *LocalLog) GetLevel() log.LogLevel {
	switch logger.Logger.Level {
	case logrus.PanicLevel:
		return log.PanicLevel
	case logrus.FatalLevel:
		return log.FatalLevel
	case logrus.ErrorLevel:
		return log.ErrorLevel
	case logrus.WarnLevel:
		return log.WarnLevel
	case logrus.InfoLevel:
		return log.InfoLevel
	case logrus.DebugLevel:
		return log.DebugLevel
	case logrus.TraceLevel:
		return log.TraceLevel
	default:
		return log.InfoLevel
	}
}

func (logger *LocalLog) SetLevel(loglevel LogLevel) {
	var LLevel logrus.Level
	switch loglevel {
	case log.PanicLevel:
		LLevel = logrus.PanicLevel
	case log.FatalLevel:
		LLevel = logrus.FatalLevel
	case log.ErrorLevel:
		LLevel = logrus.ErrorLevel
	case log.WarnLevel:
		LLevel = logrus.WarnLevel
	case log.InfoLevel:
		LLevel = logrus.InfoLevel
	case log.DebugLevel:
		LLevel = logrus.DebugLevel
	case log.TraceLevel:
		LLevel = logrus.TraceLevel
	default:
		LLevel = logrus.InfoLevel
	}

	alllogfile := filepath.Join(logger.ALL_LogfolderABS, "all_log.txt")
	errlogfile := filepath.Join(logger.ERR_LogfolderABS, "err_log.txt")

	rotateFileHook_ALL := newRotateFileHook(rotateFileConfig{
		Filename:   alllogfile,
		MaxSize:    logger.MaxSize, // megabytes
		MaxBackups: logger.MaxBackups,
		MaxAge:     logger.MaxAge, //days
		Level:      LLevel,
		Formatter: UTCFormatter{&nested.Formatter{
			NoColors:        true,
			HideKeys:        false,
			TimestampFormat: "2006-01-02 15:04:05",
		}},
	})

	rotateFileHook_ERR := newRotateFileHook(rotateFileConfig{
		Filename:   errlogfile,
		MaxSize:    logger.MaxSize, // megabytes
		MaxBackups: logger.MaxBackups,
		MaxAge:     logger.MaxAge, //days
		Level:      logrus.ErrorLevel,
		Formatter: UTCFormatter{&nested.Formatter{
			NoColors:        true,
			HideKeys:        false,
			TimestampFormat: "2006-01-02 15:04:05",
		}},
	})

	logger.SetFormatter(UTCFormatter{&nested.Formatter{
		HideKeys:        false,
		TimestampFormat: "2006-01-02 15:04:05",
	}})

	/////set hooks
	logger.Logger.SetLevel(logrus.Level(loglevel))
	logger.ReplaceHooks(make(logrus.LevelHooks))
	logger.AddHook(rotateFileHook_ALL)
	logger.AddHook(rotateFileHook_ERR)
}

// Default is info level
func New(logsAbsFolder string, fileMaxSizeMBytes int, MaxBackupsFiles int, MaxAgeDays int) (log.Logger, error) {

	logger := logrus.New()

	logsAllAbsFolder = filepath.Join(logsAbsFolder, "all")
	logsErrorAbsFolder = filepath.Join(logsAbsFolder, "error")
	//make sure the logs folder exist otherwise create dir
	err := os.MkdirAll(logsAllAbsFolder, 0777)
	if err != nil {
		return nil, err
	}
	err = os.MkdirAll(logsErrorAbsFolder, 0777)
	if err != nil {
		return nil, err
	}

	//default info level//
	LocalLogPointer := &LocalLog{logger, logsAllAbsFolder, logsErrorAbsFolder,
		fileMaxSizeMBytes, MaxBackupsFiles, MaxAgeDays}
	LocalLogPointer.SetLevel(log.InfoLevel)
	return LocalLogPointer, nil
}

func (logger *LocalLog) getLogFilesList(log_folder string) ([]string, error) {

	var result []string
	files, err := ioutil.ReadDir(log_folder)
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return result, nil
	}

	for i := len(files) - 1; i >= 0; i-- {
		result = append(result, files[i].Name())
	}

	return result, nil
}

func (logger *LocalLog) GetLastN(lineCount int64, levels []LogLevel) ([]string, error) {
	var alllogfiles []string
	var err error
	folder := logger.ALL_LogfolderABS
	alllogfiles, err = logger.getLogFilesList(folder)
	if err != nil {
		return nil, err
	}
	if len(alllogfiles) == 0 {
		return nil, errors.New("no logfile")
	}

	var Counter int64 = 0
	levelMap := map[LogLevel]struct{}{}
	for i := range levels {
		levelMap[levels[i]] = struct{}{}
	}

	resultLog := []string{}

	for i := 0; i < len(alllogfiles); i++ {
		fname := filepath.Join(folder, alllogfiles[i])

		fileContent, err := ioutil.ReadFile(fname)
		if err != nil {
			return resultLog, err
		}

		lines := splitLines(string(fileContent))

		for i := len(lines) - 1; i >= 0; i-- {

			if lines[i] == "" {
				continue
			}

			if strings.Contains(lines[i], string(log.DebugTagStr)) && isContain(levelMap, log.DebugLevel) {
				resultLog = append(resultLog, lines[i])
				Counter++
			} else if strings.Contains(lines[i], string(log.TraceTagStr)) && isContain(levelMap, log.TraceLevel) {
				resultLog = append(resultLog, lines[i])
				Counter++
			} else if strings.Contains(lines[i], string(log.InfoTagStr)) && isContain(levelMap, log.InfoLevel) {
				resultLog = append(resultLog, lines[i])
				Counter++
			} else if strings.Contains(lines[i], string(log.WarnTagStr)) && isContain(levelMap, log.WarnLevel) {
				resultLog = append(resultLog, lines[i])
				Counter++
			} else if strings.Contains(lines[i], string(log.FatalTagStr)) && isContain(levelMap, log.FatalLevel) {
				resultLog = append(resultLog, lines[i])
				Counter++
			} else if strings.Contains(lines[i], string(log.ErrorTagStr)) && isContain(levelMap, log.ErrorLevel) {
				resultLog = append(resultLog, lines[i])
				Counter++
			} else if strings.Contains(lines[i], string(log.PanicTagStr)) && isContain(levelMap, log.PanicLevel) {
				resultLog = append(resultLog, lines[i])
				Counter++
			}

			if Counter >= lineCount {
				reversArr(resultLog)
				return resultLog, nil
			}
		}
	}
	reversArr(resultLog)
	return resultLog, nil
}

func (logger *LocalLog) PrintLastN(lineCount int64, levels []LogLevel) {
	color.White("================== start ==================")

	lines, err := logger.GetLastN(lineCount, levels)
	if err != nil {
		color.Red(err.Error())
		color.White("=================== end ===================")
		return
	}

	if err != nil {
		color.Red(err.Error())
		color.White("=================== end ===================")
		return
	}

	var Counter int64 = 0
	levelMap := map[LogLevel]struct{}{}
	for i := range levels {
		levelMap[levels[i]] = struct{}{}
	}

	for i := 0; i < len(lines); i++ {

		if strings.Contains(lines[i], string(log.DebugTagStr)) && isContain(levelMap, log.DebugLevel) {
			color.White(lines[i])
			Counter++
		} else if strings.Contains(lines[i], string(log.TraceTagStr)) && isContain(levelMap, log.TraceLevel) {
			color.Cyan(lines[i])
			Counter++
		} else if strings.Contains(lines[i], string(log.InfoTagStr)) && isContain(levelMap, log.InfoLevel) {
			color.Green(lines[i])
			Counter++
		} else if strings.Contains(lines[i], string(log.WarnTagStr)) && isContain(levelMap, log.WarnLevel) {
			color.Yellow(lines[i])
			Counter++
		} else if strings.Contains(lines[i], string(log.FatalTagStr)) && isContain(levelMap, log.FatalLevel) {
			color.Red(lines[i])
			Counter++
		} else if strings.Contains(lines[i], string(log.ErrorTagStr)) && isContain(levelMap, log.ErrorLevel) {
			color.Red(lines[i])
			Counter++
		} else if strings.Contains(lines[i], string(log.PanicTagStr)) && isContain(levelMap, log.PanicLevel) {
			color.Red(lines[i])
			Counter++
		}

		if Counter >= lineCount {
			color.White("=================== end ===================")
			return
		}

	}
	color.White("=================== end ===================")
}

func splitLines(s string) []string {
	var lines []string
	sc := bufio.NewScanner(strings.NewReader(s))
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}
	return lines
}

func isContain(m map[log.LogLevel]struct{}, l log.LogLevel) bool {
	_, ok := m[l]
	return ok
}

func reversArr(arr []string) {
	length := len(arr)
	for i := 0; i < length/2; i++ {
		temp := arr[length-1-i]
		arr[length-1-i] = arr[i]
		arr[i] = temp
	}
}
