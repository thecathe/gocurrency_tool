package log

import (
	"fmt"
	"log"
	"os"
)

var enabled_general bool = true
var enabled_debug bool = true
var enabled_warning bool = true
var enabled_failure bool = true
var enabled_verbose bool = false

var total_logs int = 0

// for displaying counter
var display_counter bool = false
var counter_value int = -1

func DisplayCounter() bool {
	return display_counter
}

func SetDisplay(b bool) {
	display_counter = b
}

func SetCounter(c int) {
	counter_value = c
}

var _log = log.New(os.Stdout, "", log.Ltime|log.Lmsgprefix)

func Total() int {
	return total_logs
}

func CheckCounter(s string) string {
	if display_counter {
		s = fmt.Sprintf("| %04d | %s", counter_value, s)
	}
	return s
}

func SetLoggers(_general bool, _debug bool, _warning bool, _failure bool) {
	enabled_general = _general
	enabled_debug = _debug
	enabled_warning = _warning
	enabled_failure = _failure
}

func GeneralLog(s string, v ...interface{}) {
	if enabled_general {
		_log.Printf(CheckCounter("General: "+s), v...)
		total_logs++
	}
}

func DebugLog(s string, v ...interface{}) {
	if enabled_debug {
		_log.Printf(CheckCounter("Debug: "+s), v...)
		total_logs++
	}
}

func WarningLog(s string, v ...interface{}) {
	if enabled_warning {
		_log.Printf(CheckCounter("Warning: "+s), v...)
		total_logs++
	}
}

func FailureLog(s string, v ...interface{}) {
	if enabled_failure {
		_log.Printf(CheckCounter("Failure: "+s), v...)
		total_logs++
	}
}

// internally modifiable logs
func VerboseLog(s string, v ...interface{}) {
	if enabled_verbose {
		_log.Printf(CheckCounter("Verbose: "+s), v...)
		total_logs++
	}
}

func ExitLog(i int, s string, v ...interface{}) {
	DebugLog("Total Logs: %d\n", total_logs)
	_log.Printf("Exiting: "+s, v...)
	os.Exit(i)
}

func PanicLog(e error, s string, v ...interface{}) {
	DebugLog("Total Logs: %d\n", total_logs)
	_log.Printf("Panicking: "+s, v...)
	panic(e)
}
