package lumberjackrus

import (
	"testing"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"bytes"
	"os"
)

const expectedMsg = "Test message."
const unexpectedMsg = "Second test message"

const traceMsg = "Trace message"
const debugMsg = "Debug message"
const infoMsg = "Info message"
const errorMsg = "Error message"

const defaultFile = "/tmp/lumberjackrus_default.log"
const debugFile = "/tmp/lumberjackrus_debug.log"
const errorFile = "/tmp/lumberjackrus_error.log"

func initHookWithoutOption(formatter logrus.Formatter, minLevel logrus.Level) (*Hook, error) {
	return NewHook(
		&LogFile{
			Filename:   defaultFile,
			MaxSize:    100,
			MaxBackups: 1,
			MaxAge:     1,
			Compress:   false,
		},
		minLevel,
		formatter,
		nil,
	)
}

func initHookWithOption(formatter logrus.Formatter, minLevel logrus.Level) (*Hook, error) {

	return NewHook(
		&LogFile{
			Filename:   defaultFile,
			MaxSize:    100,
			MaxBackups: 1,
			MaxAge:     1,
			Compress:   false,
		},
		minLevel,
		formatter,
		&LogFileOpts{
			logrus.DebugLevel: &LogFile{Filename: debugFile},
			logrus.ErrorLevel: &LogFile{Filename: errorFile},
		},
	)
}

func TestDefaultLogger(t *testing.T) {
	defer func() {
		os.Remove(defaultFile)
	}()

	// Init
	minLevel := logrus.DebugLevel
	formatter := &logrus.TextFormatter{}

	log := logrus.New()
	log.SetLevel(minLevel + 1)

	hook, err := initHookWithoutOption(formatter, logrus.InfoLevel)
	if err != nil {
		t.Errorf("Unable to instantiate new hook: %s", err)
	}

	log.AddHook(hook)

	// Run
	log.Debug(unexpectedMsg)
	log.Info(expectedMsg)

	// Assertion
	assertLogs(t, defaultFile, expectedMsg, unexpectedMsg)
}

func TestOptionsLogger(t *testing.T) {
	defer func() {
		os.Remove(defaultFile)
		os.Remove(debugFile)
		os.Remove(errorFile)
	}()

	// Init
	minLevel := logrus.DebugLevel
	formatter := &logrus.TextFormatter{}

	log := logrus.New()
	log.SetLevel(minLevel + 1)

	hook, err := initHookWithOption(formatter, logrus.DebugLevel)
	if err != nil {
		t.Errorf("Unable to instantiate new hook: %s", err)
	}

	log.AddHook(hook)

	// Run
	log.Trace(traceMsg)
	log.Debug(debugMsg)
	log.Info(infoMsg)
	log.Error(errorMsg)

	// Assertion
	assertLogs(t, defaultFile, infoMsg, traceMsg, debugMsg, errorMsg)
	assertLogs(t, debugFile, debugMsg, traceMsg, infoMsg, errorMsg)
	assertLogs(t, errorFile, errorMsg, traceMsg, debugMsg, infoMsg)
}

func assertLogs(t *testing.T, filename, expectedMessage string, unexpectedMessages ...string) {
	logs, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Errorf("Error while reading from log file: %s", err)
		return
	}
	assertLogsContains(t, logs, filename, expectedMessage)

	for _, unexpected := range unexpectedMessages {
		assertLogsNotContains(t, logs, filename, unexpected)
	}
}

func assertLogsContains(t *testing.T, logs []byte, filename, expected string) {
	if !bytes.Contains(logs, []byte("msg=\""+expected+"\"")) {
		t.Errorf("%s doesn't contain the expected message '%s'", filename, expected)
	}
}

func assertLogsNotContains(t *testing.T, logs []byte, filename, unexpected string) {
	if bytes.Contains(logs, []byte("msg=\""+unexpected+"\"")) {
		t.Errorf("%s contains the unexpected message '%s'", filename, unexpected)
	}
}
