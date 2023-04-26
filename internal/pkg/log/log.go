// Package log provides logging functions
// wraps fmt.Println with options and colors
package log

import (
	"fmt"
	"io"

	"github.com/fatih/color"
)

var (
	// DiagnosticWriter controls logs output
	DiagnosticWriter = color.Error
	// OutputWriter     = color.Output
)

// Colored string formatting functions.
var (
	successSprintf = color.HiGreenString
	errorSprintf   = color.HiRedString
	warningSprintf = color.YellowString
	debugSprintf   = color.New(color.Faint).Sprintf
)

// Log message prefixes.
var (
	successPrefix = "✔"
	errorPrefix   = "error:"
)

// Info writes the message to standard error with the default color.
func Info(args ...interface{}) {
	info(DiagnosticWriter, args...)
}

func info(w io.Writer, args ...interface{}) {
	fmt.Fprintln(w, args...)
}

// Infof formats according to the specifier, and writes to standard error with the default color.
func Infof(format string, args ...interface{}) {
	infof(DiagnosticWriter, format, args...)
}

func infof(w io.Writer, format string, args ...interface{}) {
	fmt.Fprintf(w, format, args...)
}

// Error writes the message to standard error with the default color.
func Error(args ...interface{}) {
	error(DiagnosticWriter, args...)
}

//revive:disable
func error(w io.Writer, args ...interface{}) {
	msg := fmt.Sprintf("%s %s", errorSprintf(errorPrefix), fmt.Sprint(args...))
	fmt.Fprintln(w, msg)
}

// Success prefixes the message with a green "✔ Success!", and writes to standard error with a new line.
func Success(args ...interface{}) {
	success(DiagnosticWriter, args...)
}

func success(w io.Writer, args ...interface{}) {
	msg := fmt.Sprintf("%s %s", successSprintf(successPrefix), fmt.Sprint(args...))
	fmt.Fprintln(w, msg)
}

// Debug writes the message to standard error in grey and with a new line.
func Debug(args ...interface{}) {
	debug(DiagnosticWriter, args...)
}

func debug(w io.Writer, args ...interface{}) {
	fmt.Fprintln(w, debugSprintf(fmt.Sprint(args...)))
}

// Warning prefixes the message with a "Note:", colors the *entire* message in yellow, writes to standard error with a new line.
func Warning(args ...interface{}) {
	warning(DiagnosticWriter, args...)
}
func warning(w io.Writer, args ...interface{}) {
	msg := fmt.Sprint(args...)
	fmt.Fprintln(w, warningSprintf(fmt.Sprintf("%s %s", "warning: ", msg)))
}

// Debugf formats according to the specifier, colors the message in grey, and writes to standard error.
func Debugf(format string, args ...interface{}) {
	debugf(DiagnosticWriter, format, args...)
}

func debugf(w io.Writer, format string, args ...interface{}) {
	fmt.Fprint(w, debugSprintf(format, args...))
}
