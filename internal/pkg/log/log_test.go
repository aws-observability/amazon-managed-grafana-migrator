package log

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const verbose = false

func TestInfo(t *testing.T) {

	b := &strings.Builder{}
	DiagnosticWriter = b

	// WHEN
	Info("hello", "world")

	// THEN
	require.Equal(t, b.String(), "hello world\n")
}

func TestError(t *testing.T) {
	// GIVEN
	b := &strings.Builder{}
	DiagnosticWriter = b

	// WHEN
	Error("hello", " world")

	// THEN
	require.Contains(t, b.String(), fmt.Sprintf("%s hello world\n", errorPrefix))
}

func TestInfof(t *testing.T) {
	// GIVEN
	b := &strings.Builder{}
	DiagnosticWriter = b

	// WHEN
	Infof("%s %s\n", "hello", "world")

	// THEN
	require.Equal(t, "hello world\n", b.String())
}

func TestSuccess(t *testing.T) {
	// GIVEN
	b := &strings.Builder{}
	DiagnosticWriter = b

	// WHEN
	Success("hello", " world")

	// THEN
	require.Equal(t, b.String(), fmt.Sprintf("%s hello world\n", successPrefix))
}

func TestDebug(t *testing.T) {
	// GIVEN
	b := &strings.Builder{}
	DiagnosticWriter = b

	// WHEN
	Debug(verbose, "hello", " world")

	// THEN
	require.Contains(t, b.String(), "")

	// WHEN
	Debug(true, "hello", " world")

	// THEN
	require.Contains(t, b.String(), "hello world\n")
}

func TestWarning(t *testing.T) {
	// GIVEN
	b := &strings.Builder{}
	DiagnosticWriter = b

	// WHEN
	Warning("hello", " world")

	// THEN
	require.Contains(t, b.String(), "warning:")
	require.Contains(t, b.String(), "hello world\n")
}

func TestDebugf(t *testing.T) {
	// GIVEN
	b := &strings.Builder{}
	DiagnosticWriter = b

	// WHEN
	Debugf(verbose, "%s %s\n", "hello", "world")

	// THEN
	require.Contains(t, b.String(), "")

	// WHEN
	Debugf(true, "%s %s\n", "hello", "world")

	// THEN
	require.Contains(t, b.String(), "hello world\n")
}
