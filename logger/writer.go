package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var std = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)

type Writer struct {
	Level Level
	Std   bool

	prefix string
	folder string

	logger *log.Logger
	year   int
	month  time.Month
	day    int
	file   *os.File
	mu     sync.Mutex
}

func (s *Writer) Init(level string, prefix, folder string) error {
	s.Level.Parse(level)
	s.prefix = prefix
	s.folder = folder

	return nil
}

func (s *Writer) Close() {
	if s.file != nil {
		s.file.Close()
		s.file = nil
	}
}

func (s *Writer) Error(v ...interface{}) string {
	return s.output(LevelError, fmt.Sprint(v...))
}

func (s *Writer) Warning(v ...interface{}) string {
	return s.output(LevelWarning, fmt.Sprint(v...))
}

func (s *Writer) Info(v ...interface{}) string {
	return s.output(LevelInfo, fmt.Sprint(v...))
}

func (s *Writer) Trace(v ...interface{}) string {
	return s.output(LevelTrace, fmt.Sprint(v...))
}

func (s *Writer) Debug(v ...interface{}) string {
	return s.output(LevelDebug, fmt.Sprint(v...))
}

func (s *Writer) getLogger() *log.Logger {
	if s.folder != "" {
		now := time.Now()
		if s.year != now.Year() || s.month != now.Month() || s.day != now.Day() || s.file == nil {
			s.mu.Lock()
			defer s.mu.Unlock()

			s.year = now.Year()
			s.month = now.Month()
			s.day = now.Day()
			if s.logger != nil {
				s.logger = nil
			}
			if s.file != nil {
				s.file.Close()
				s.file = nil
			}

			os.MkdirAll(s.folder, 0777)
			prefix := ""
			if s.prefix != "" {
				prefix = fmt.Sprintf("%s_", s.prefix)
			}
			fileName := fmt.Sprint(prefix, now.Year(), "-", int(now.Month()), "-", now.Day(), ".log")
			filePath := filepath.Join(s.folder, fileName)
			file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
			if err != nil {
				return nil
			}
			s.file = file
			s.logger = log.New(s.file, "", log.Ldate|log.Ltime|log.Lshortfile)
		}
	} else {
		if s.logger == nil {
			prefix := ""
			if s.prefix != "" {
				prefix = fmt.Sprintf("%s ", s.prefix)
			}
			s.logger = log.New(os.Stderr, prefix, log.Ldate|log.Ltime|log.Lshortfile)
		}
	}

	return s.logger
}

func (s *Writer) output(l Level, m string) string {
	str := fmt.Sprintf("%s; %s", levelText[l], m)

	if s.Level&l != 0 {
		if s.Std && s.file != nil {
			std.Output(4, fmt.Sprintln(str))
		}

		logger := s.getLogger()
		if logger != nil {
			logger.Output(4, fmt.Sprintln(str))
		}
	}

	return str
}
