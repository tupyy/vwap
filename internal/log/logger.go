// This package provide a simple wrapper around standard logger in order to provide level and method name functionality.
package log

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
)

type Level int

func (ll Level) String() string {
	switch ll {
	case Trace:
		return "TRACE"
	case Debug:
		return "DEBUG"
	case Info:
		return "INFO"
	case Warning:
		return "WARNING"
	case Error:
		return "ERROR"
	default:
		return "unknown level"
	}
}

const (
	Trace Level = iota
	Debug
	Info
	Warning
	Error
)

type Logger struct {
	level      Level
	methodName string
	logger     *log.Logger
}

var logger *Logger

func newLogger(logLevel Level) *Logger {
	return &Logger{
		level:  logLevel,
		logger: log.New(os.Stdout, "", log.LstdFlags),
	}
}

func SetLogLevel(level Level) {
	if logger == nil {
		logger = newLogger(level)
	} else {
		logger.level = level
	}
}

func GetLogger() *Logger {
	if logger == nil {
		logger = newLogger(Info)
	}

	logger.methodName = getMethodName()

	return logger
}

func (l *Logger) Tracef(format string, v ...interface{}) {
	l.outputf(Trace, format, v...)
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	l.outputf(Debug, format, v...)
}

func (l *Logger) Infof(format string, v ...interface{}) {
	l.outputf(Info, format, v...)
}

func (l *Logger) Warningf(format string, v ...interface{}) {
	l.outputf(Warning, format, v...)
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	l.outputf(Error, format, v...)
}

func (l *Logger) outputf(logLevel Level, format string, v ...interface{}) {
	if l.level <= logLevel {
		l.logger.Printf("[%s] Method[%s] %s", logLevel.String(), l.methodName, fmt.Sprintf(format, v...))
	}
}

// getMethodName get the calling method from the stack.
func getMethodName() string {
	// 4 as stack depth should be enough to get the real caller. (2 should be enough)
	stack := make([]uintptr, 4)
	depth := runtime.Callers(3, stack) // Can skip the first 3 as it's Callers < getMethodName < Get(*)Logger

	if depth < 1 {
		return ""
	}

	frames := runtime.CallersFrames(stack)

	for f, hasNext := frames.Next(); hasNext; {
		tmp := strings.Split(f.Function, "/")
		if len(tmp) == 0 {
			continue
		}

		shortName := tmp[len(tmp)-1]

		if !strings.HasPrefix(shortName, "log.") {
			return shortName
		}
	}

	return ""
}
