package v1

import (
	"strings"

	"github.com/google/uuid"
)

// generateUserID creates a unique user ID
func generateUserID() string {
	return uuid.New().String()
}

var (
	logLevelStringsDebug = []string{"debug", "dbg", "debugging"}
	logLevelStringsInfo  = []string{"info", "inf", "information"}
	logLevelStringsWarn  = []string{"warn", "wrn", "warning"}
	logLevelStringsError = []string{"error", "err"}
	logLevelStringsFatal = []string{"fatal", "critical"}
)

var logLevelEmoji = map[string]string{
	logLevelStringsDebug[0]: "🐛",
	logLevelStringsInfo[0]:  "💬",
	logLevelStringsWarn[0]:  "⚠️",
	logLevelStringsError[0]: "❌",
	logLevelStringsFatal[0]: "💣",
}

func getLogLevelEmoji(level string) string {
	level = strings.ToLower(level)

	logLevelStrings := [][]string{
		logLevelStringsDebug,
		logLevelStringsInfo,
		logLevelStringsWarn,
		logLevelStringsError,
		logLevelStringsFatal,
	}

	for _, llstrings := range logLevelStrings {
		for _, llstring := range llstrings {
			if llstring == level {
				return logLevelEmoji[llstrings[0]]
			}
		}
	}

	return ""
}
