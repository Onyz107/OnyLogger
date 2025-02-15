package logger

import (
	"testing"
	"time"
)

func TestSpinner(t *testing.T) {
	spinner := NewSpinner("Building go")
	spinner.Start()

	time.Sleep(5 * time.Second)

	spinner.Stop()
}
