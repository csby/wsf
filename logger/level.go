package logger

import "strings"

const (
	LevelError = 1 << iota
	LevelWarning
	LevelInfo
	LevelTrace
	LevelDebug
	LevelAll = LevelError | LevelWarning | LevelInfo | LevelTrace | LevelDebug
)

type Level int

var levelText = map[Level]string{
	LevelError:   "error",
	LevelWarning: "warning",
	LevelInfo:    "info",
	LevelTrace:   "trace",
	LevelDebug:   "debug",
}

func (s Level) String() string {
	levels := make([]string, 0)

	if s&LevelError != 0 {
		levels = append(levels, levelText[LevelError])
	}
	if s&LevelWarning != 0 {
		levels = append(levels, levelText[LevelWarning])
	}
	if s&LevelInfo != 0 {
		levels = append(levels, levelText[LevelInfo])
	}
	if s&LevelTrace != 0 {
		levels = append(levels, levelText[LevelTrace])
	}
	if s&LevelDebug != 0 {
		levels = append(levels, levelText[LevelDebug])
	}

	if len(levels) > 0 {
		return strings.Join(levels, "|")
	}

	return ""
}

func (s *Level) Parse(level string) {
	*s = 0
	levels := strings.Split(strings.ToLower(level), "|")
	lenth := len(levels)
	for index := 0; index < lenth; index++ {
		if strings.TrimSpace(levels[index]) == "error" {
			*s |= LevelError
		}
		if strings.TrimSpace(levels[index]) == "warning" {
			*s |= LevelWarning
		}
		if strings.TrimSpace(levels[index]) == "info" {
			*s |= LevelInfo
		}
		if strings.TrimSpace(levels[index]) == "trace" {
			*s |= LevelTrace
		}
		if strings.TrimSpace(levels[index]) == "debug" {
			*s |= LevelDebug
		}
	}
}
