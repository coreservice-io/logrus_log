package Logrus

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
	"github.com/universe-30/Logrus/nested"
	"github.com/universe-30/ULog"
)

var logsAllAbsFolder string
var logsErrorAbsFolder string

type Fields = logrus.Fields
type LogLevel = ULog.LogLevel

type LocalLog struct {
	*logrus.Logger
	ALL_LogfolderABS string
	ERR_LogfolderABS string
	MaxSize          int
	MaxBackups       int
	MaxAge           int
}

func (logger *LocalLog) SetLevel(loglevel LogLevel) {

	var LLevel logrus.Level

	switch loglevel {
	case ULog.PanicLevel:
		LLevel = logrus.PanicLevel
	case ULog.FatalLevel:
		LLevel = logrus.FatalLevel
	case ULog.ErrorLevel:
		LLevel = logrus.ErrorLevel
	case ULog.WarnLevel:
		LLevel = logrus.WarnLevel
	case ULog.InfoLevel:
		LLevel = logrus.InfoLevel
	case ULog.DebugLevel:
		LLevel = logrus.DebugLevel
	case ULog.TraceLevel:
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
		//NoColors:        !ShowColor,
	}})

	/////set hooks
	logger.Logger.SetLevel(logrus.Level(loglevel))
	logger.ReplaceHooks(make(logrus.LevelHooks))
	logger.AddHook(rotateFileHook_ALL)
	logger.AddHook(rotateFileHook_ERR)
}

// Default is info level
func New(logsAbsFolder string, fileMaxSizeMBytes int, MaxBackupsFiles int, MaxAgeDays int) (*LocalLog, error) {

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
	LocalLogPointer.SetLevel(ULog.InfoLevel)
	return LocalLogPointer, nil
}

func (logger *LocalLog) GetLogFilesList(log_folder string) ([]string, error) {

	var result []string
	files, err := ioutil.ReadDir(log_folder)
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		result = append(result, f.Name())
	}
	return result, nil
}

func (logger *LocalLog) PrintLastN_ErrLogs(lastN int) {
	logger.printLastNLogs("error", lastN)
}

func (logger *LocalLog) PrintLastN_AllLogs(lastN int) {
	logger.printLastNLogs("all", lastN)
}

func (logger *LocalLog) printLastNLogs(type_ string, lastN int) {

	color.White("================== start ==================")

	var alllogfiles []string
	var err error
	var folder string
	if type_ == "error" {
		folder = logger.ERR_LogfolderABS
	} else {
		folder = logger.ALL_LogfolderABS
	}
	alllogfiles, err = logger.GetLogFilesList(folder)

	if err != nil {
		color.Red(err.Error())
		color.White("================== end   ==================")
		return
	}
	if len(alllogfiles) == 0 {
		color.Red("no logfile")
		color.White("================== end   ==================")
		return
	}

	Counter := 0

	for i := 0; i < len(alllogfiles); i++ {
		fname := filepath.Join(folder, alllogfiles[i])

		var cmd *exec.Cmd

		if runtime.GOOS == "windows" {
			cmd = exec.Command("powershell", "-nologo", "-noprofile")
			stdin, err := cmd.StdinPipe()
			if err != nil {
				color.Red(err.Error())
				color.Red("log view not supported , please directly check logfile :" + fname)
				color.White("================== end   ==================")
				return
			}
			go func() {
				defer stdin.Close()
				fmt.Fprintln(stdin, "Get-Content "+fname+" | Select-Object -last "+strconv.Itoa(lastN))
			}()
		} else {
			cmd = exec.Command("tail", "-n", strconv.Itoa(lastN), fname)
		}

		stdout, err := cmd.Output()
		if err != nil {
			color.Red(err.Error())
			color.White("================== end   ==================")
			return
		}
		lines := splitLines(string(stdout))
		for i := 0; i < len(lines); i++ {

			if strings.Contains(lines[i], string(ULog.DebugTagStr)) {
				color.White(lines[i])
			} else if strings.Contains(lines[i], string(ULog.TraceTagStr)) {
				color.Cyan(lines[i])
			} else if strings.Contains(lines[i], string(ULog.InfoTagStr)) {
				color.Green(lines[i])
			} else if strings.Contains(lines[i], string(ULog.WarnTagStr)) {
				color.Yellow(lines[i])
			} else if strings.Contains(lines[i], string(ULog.FatalTagStr)) ||
				strings.Contains(lines[i], string(ULog.ErrorTagStr)) ||
				strings.Contains(lines[i], string(ULog.PanicTagStr)) {
				color.Red(lines[i])
			} else {
				color.White(lines[i])
			}

			Counter++
			if Counter >= lastN {
				color.White("================== end   ==================")
				return
			}
		}

	}
	color.White("================== end   ==================")
}

func splitLines(s string) []string {
	var lines []string
	sc := bufio.NewScanner(strings.NewReader(s))
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}
	return lines
}
