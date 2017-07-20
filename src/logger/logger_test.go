package logger

import (
	"testing"
)

var logger *Logger = NewLoggerWithFile(LEVEL_INFO, "log2filetest.log")

func TestDEBUG(t *testing.T) {
	logger.DEBUG("level debug")
}

func TestInfo(t *testing.T) {
	logger.INFO("level info")
}

func TestFATAL(t *testing.T) {
	logger.FATAL("level fatal")
}
