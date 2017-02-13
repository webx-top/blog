// Copyright 2015 Qiang Xue. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package log implements logging with severity levels and message categories.
package log

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// RFC5424 log message levels.
const (
	LevelFatal Level = iota
	LevelError
	LevelWarn
	LevelInfo
	LevelDebug
)

const (
	ActionNothing Action = iota
	ActionPanic
	ActionExit
)

// Level describes the level of a log message.
type Level int
type Action int

// LevelNames maps log levels to names
var LevelNames = map[Level]string{
	LevelDebug: "Debug",
	LevelInfo:  "Info",
	LevelWarn:  "Warn",
	LevelError: "Error",
	LevelFatal: "Fatal",
}

var Levels = map[string]Level{
	"Debug": LevelDebug,
	"Info":  LevelInfo,
	"Warn":  LevelWarn,
	"Error": LevelError,
	"Fatal": LevelFatal,
}

func GetLevel(level string) (Level, bool) {
	level = strings.Title(level)
	l, y := Levels[level]
	return l, y
}

// String returns the string representation of the log level
func (l Level) String() string {
	if name, ok := LevelNames[l]; ok {
		return name
	}
	return "Unknown"
}

type LoggerWriter struct {
	Level Level
	*Logger
}

func (l *LoggerWriter) Write(p []byte) (n int, err error) {
	var s string
	n = len(p)
	if p[n-1] == '\n' {
		s = string(p[0 : n-1])
	} else {
		s = string(p)
	}
	l.Logger.newEntry(l.Level, s)
	return
}

// Entry represents a log entry.
type Entry struct {
	Level     Level
	Category  string
	Message   string
	Time      time.Time
	CallStack string

	FormattedMessage string
}

// String returns the string representation of the log entry
func (e *Entry) String() string {
	return e.FormattedMessage
}

// Target represents a target where the logger can send log messages to for further processing.
type Target interface {
	// Open prepares the target for processing log messages.
	// Open will be invoked when Logger.Open() is called.
	// If an error is returned, the target will be removed from the logger.
	// errWriter should be used to write errors found while processing log messages.
	Open(errWriter io.Writer) error
	// Process processes an incoming log message.
	Process(*Entry)
	// Close closes a target.
	// Close is called when Logger.Close() is called, which gives each target
	// a chance to flush the logged messages to their destination storage.
	Close()
	SetLevel(interface{})
	SetLevels(...Level)
}

// coreLogger maintains the log messages in a channel and sends them to various targets.
type coreLogger struct {
	lock        sync.Mutex
	open        bool        // whether the logger is open
	entries     chan *Entry // log entries
	goroutines  int32
	fatalAction Action

	ErrorWriter     io.Writer // the writer used to write errors caused by log targets
	BufferSize      int       // the size of the channel storing log entries
	CallStackDepth  int       // the number of call stack frames to be logged for each message. 0 means do not log any call stack frame.
	CallStackFilter string    // a substring that a call stack frame file path should contain in order for the frame to be counted
	MaxLevel        Level     // the maximum level of messages to be logged
	Targets         []Target  // targets for sending log messages to
	SyncMode        bool      // Whether the use of non-asynchronous mode （是否使用非异步模式）
	MaxGoroutines   int32     // Max Goroutine
	AddSpace        bool      // Add a space between two arguments.
}

// Formatter formats a log message into an appropriate string.
type Formatter func(*Logger, *Entry) string

// Logger records log messages and dispatches them to various targets for further processing.
type Logger struct {
	*coreLogger
	Category   string    // the category associated with this logger
	Formatter  Formatter // message formatter
	categories map[string]*Logger
}

// NewLogger creates a root logger.
// The new logger takes these default options:
// ErrorWriter: os.Stderr, BufferSize: 1024, MaxLevel: LevelDebug,
// Category: app, Formatter: DefaultFormatter
func NewLogger(args ...string) *Logger {
	logger := &coreLogger{
		ErrorWriter:   os.Stderr,
		BufferSize:    1024,
		MaxLevel:      LevelDebug,
		Targets:       make([]Target, 0),
		MaxGoroutines: 100000,
	}
	category := `app`
	if len(args) > 0 {
		category = args[0]
	}
	logger.Targets = append(logger.Targets, NewConsoleTarget())
	logger.Open()
	return &Logger{
		coreLogger: logger,
		Category:   category,
		Formatter:  NormalFormatter,
		categories: make(map[string]*Logger),
	}
}

func New(args ...string) *Logger {
	return NewLogger(args...)
}

// GetLogger creates a logger with the specified category and log formatter.
// Messages logged through this logger will carry the same category name.
// The formatter, if not specified, will inherit from the calling logger.
// It will be used to format all messages logged through this logger.
func (l *Logger) GetLogger(category string, formatter ...Formatter) *Logger {
	logger, ok := l.categories[category]
	if !ok {
		logger = &Logger{
			coreLogger: l.coreLogger,
			Category:   category,
		}
		if len(formatter) > 0 {
			logger.Formatter = formatter[0]
		} else {
			logger.Formatter = l.Formatter
		}
		l.categories[category] = logger
	} else {
		if len(formatter) > 0 {
			logger.Formatter = formatter[0]
		}
	}
	return logger
}

func (l *Logger) Sync(args ...bool) {
	if len(args) < 1 {
		l.SyncMode = true
		return
	}
	l.SyncMode = args[0]
}

func (l *Logger) SetTarget(targets ...Target) {
	l.Close()
	if len(targets) > 0 {
		l.Targets = targets
		l.Open()
	} else {
		l.Targets = []Target{}
	}
}

func (l *Logger) SetFatalAction(action Action) {
	l.fatalAction = action
}

func (l *Logger) AddTarget(targets ...Target) {
	l.Close()
	l.Targets = append(l.Targets, targets...)
	l.Open()
}

func (l *Logger) SetLevel(level string) {
	if le, ok := GetLevel(level); ok {
		l.MaxLevel = le
	}
}

func (l *Logger) Fatalf(format string, a ...interface{}) {
	l.Logf(LevelFatal, format, a...)
}

// Errorf logs a message indicating an error condition.
// This method takes one or multiple parameters. If a single parameter
// is provided, it will be treated as the log message. If multiple parameters
// are provided, they will be passed to fmt.Sprintf() to generate the log message.
func (l *Logger) Errorf(format string, a ...interface{}) {
	l.Logf(LevelError, format, a...)
}

// Warnf logs a message indicating a warning condition.
// Please refer to Error() for how to use this method.
func (l *Logger) Warnf(format string, a ...interface{}) {
	l.Logf(LevelWarn, format, a...)
}

// Infof logs a message for informational purpose.
// Please refer to Error() for how to use this method.
func (l *Logger) Infof(format string, a ...interface{}) {
	l.Logf(LevelInfo, format, a...)
}

// Debugf logs a message for debugging purpose.
// Please refer to Error() for how to use this method.
func (l *Logger) Debugf(format string, a ...interface{}) {
	l.Logf(LevelDebug, format, a...)
}

// Logf logs a message of a specified severity level.
func (l *Logger) Logf(level Level, format string, a ...interface{}) {
	if level > l.MaxLevel || !l.open {
		return
	}
	message := format
	if len(a) > 0 {
		message = fmt.Sprintf(format, a...)
	}
	l.newEntry(level, message)
}

func (l *Logger) Writer(level Level) io.Writer {
	return &LoggerWriter{
		Level:  level,
		Logger: l,
	}
}

func (l *Logger) Fatal(a ...interface{}) {
	l.Log(LevelFatal, a...)
}

// Error logs a message indicating an error condition.
// This method takes one or multiple parameters. If a single parameter
// is provided, it will be treated as the log message. If multiple parameters
// are provided, they will be passed to fmt.Sprintf() to generate the log message.
func (l *Logger) Error(a ...interface{}) {
	l.Log(LevelError, a...)
}

// Warn logs a message indicating a warning condition.
// Please refer to Error() for how to use this method.
func (l *Logger) Warn(a ...interface{}) {
	l.Log(LevelWarn, a...)
}

// Info logs a message for informational purpose.
// Please refer to Error() for how to use this method.
func (l *Logger) Info(a ...interface{}) {
	l.Log(LevelInfo, a...)
}

// Debug logs a message for debugging purpose.
// Please refer to Error() for how to use this method.
func (l *Logger) Debug(a ...interface{}) {
	l.Log(LevelDebug, a...)
}

// Log logs a message of a specified severity level.
func (l *Logger) Log(level Level, a ...interface{}) {
	if level > l.MaxLevel || !l.open {
		return
	}
	var message string
	if l.AddSpace {
		message = fmt.Sprintln(a...)
		message = message[:len(message)-1]
	} else {
		message = fmt.Sprint(a...)
	}
	l.newEntry(level, message)
}

// Log logs a message of a specified severity level.
func (l *Logger) newEntry(level Level, message string) {
	if level == LevelFatal {
		l.newFatalEntry(level, message)
		return
	}
	entry := &Entry{
		Category: l.Category,
		Level:    level,
		Message:  message,
		Time:     time.Now(),
	}
	if l.CallStackDepth > 0 {
		entry.CallStack = GetCallStack(3, l.CallStackDepth, l.CallStackFilter)
	}
	entry.FormattedMessage = l.Formatter(l, entry)
	if l.SyncMode {
		l.syncProcess(entry)
	} else {
		send := func() {
			atomic.AddInt32(&l.goroutines, 1)
			l.entries <- entry
		}

		if atomic.LoadInt32(&l.goroutines) < l.MaxGoroutines {
			go send()
		} else {
			send()
		}
	}
}

func (l *Logger) newFatalEntry(level Level, message string) {
	entry := &Entry{
		Category: l.Category,
		Level:    level,
		Message:  message,
		Time:     time.Now(),
	}
	stackDepth := l.CallStackDepth
	if stackDepth == 0 {
		stackDepth = 20
	}
	entry.CallStack = GetCallStack(3, stackDepth, l.CallStackFilter)
	entry.FormattedMessage = l.Formatter(l, entry)
	if l.SyncMode {
		l.syncProcess(entry)
	} else {
		atomic.AddInt32(&l.goroutines, 1)
		l.entries <- entry
	}

	for {
		goroutines := atomic.LoadInt32(&l.goroutines)
		//fmt.Println(`waiting ...`, goroutines)
		if goroutines <= 0 {
			switch l.fatalAction {
			case ActionPanic:
				panic(entry.FormattedMessage)
			case ActionExit:
				entry := &Entry{
					Category: l.Category,
					Level:    LevelWarn,
					Message:  message + `[Forced to exit]`,
					Time:     time.Now(),
				}
				entry.FormattedMessage = l.Formatter(l, entry)
				l.syncProcess(entry)
				os.Exit(-1)
			}
			break
		}
		time.Sleep(time.Duration(goroutines) * time.Microsecond)
	}
}

// Open prepares the logger and the targets for logging purpose.
// Open must be called before any message can be logged.
func (l *coreLogger) Open() error {
	l.lock.Lock()
	defer l.lock.Unlock()

	if l.open {
		return nil
	}

	if l.ErrorWriter == nil {
		return errors.New("Logger.ErrorWriter must be set.")
	}
	if l.BufferSize < 0 {
		return errors.New("Logger.BufferSize must be no less than 0.")
	}
	if l.CallStackDepth < 0 {
		return errors.New("Logger.CallStackDepth must be no less than 0.")
	}

	l.entries = make(chan *Entry, l.BufferSize)
	var targets []Target
	for _, target := range l.Targets {
		if err := target.Open(l.ErrorWriter); err != nil {
			fmt.Fprintf(l.ErrorWriter, "Failed to open target: %v\n", err)
		} else {
			targets = append(targets, target)
		}
	}
	l.Targets = targets

	go l.process()

	l.open = true

	return nil
}

// process sends the messages to targets for processing.
func (l *coreLogger) process() {
	for {
		entry := <-l.entries
		for _, target := range l.Targets {
			target.Process(entry)
		}
		atomic.AddInt32(&l.goroutines, -1)

		if entry == nil {
			break
		}
	}
}

func (l *coreLogger) syncProcess(entry *Entry) {
	if entry == nil {
		return
	}
	for _, target := range l.Targets {
		target.Process(entry)
	}
}

// Close closes the logger and the targets.
// Existing messages will be processed before the targets are closed.
// New incoming messages will be discarded after calling this method.
func (l *coreLogger) Close() {
	if !l.open {
		return
	}
	l.open = false
	// use a nil entry to signal the close of logger
	l.entries <- nil
	for _, target := range l.Targets {
		target.Close()
	}
}

// DefaultFormatter is the default formatter used to format every log message.
func DefaultFormatter(l *Logger, e *Entry) string {
	return e.Time.Format(time.RFC3339) + "|" + e.Level.String() + "|" + e.Category + "|" + e.Message + e.CallStack
}

func NormalFormatter(l *Logger, e *Entry) string {
	return e.Time.Format(`2006-01-02 15:04:05`) + "|" + e.Level.String() + "|" + e.Category + "|" + e.Message + e.CallStack
}

type JSONL struct {
	Time      string          `bson:"time" json:"time"`
	Level     string          `bson:"level" json:"level"`
	Category  string          `bson:"category" json:"category"`
	Message   json.RawMessage `bson:"message" json:"message"`
	CallStack string          `bson:"callStack" json:"callStack"`
}

func JSONFormatter(l *Logger, e *Entry) string {
	jsonl := &JSONL{
		Time:      e.Time.Format(`2006-01-02 15:04:05`),
		Level:     e.Level.String(),
		Category:  e.Category,
		Message:   []byte(`"` + e.Message + `"`),
		CallStack: e.CallStack,
	}
	if len(e.Message) > 0 {
		switch e.Message[0] {
		case '{', '[', '"':
			jsonl.Message = []byte(e.Message)
		}
	}
	b, err := json.Marshal(jsonl)
	if err != nil {
		fmt.Println(err.Error())
	}
	return string(b)
}

// GetCallStack returns the current call stack information as a string.
// The skip parameter specifies how many top frames should be skipped, while
// the frames parameter specifies at most how many frames should be returned.
func GetCallStack(skip int, frames int, filter string) string {
	buf := new(bytes.Buffer)
	for i, count := skip, 0; count < frames; i++ {
		_, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		if filter == "" || strings.Contains(file, filter) {
			fmt.Fprintf(buf, "\n%s:%d", file, line)
			count++
		}
	}
	return buf.String()
}
