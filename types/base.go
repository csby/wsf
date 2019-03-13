package types

type Base struct {
	log Log
}

func (s *Base) SetLog(log Log) {
	s.log = log
}

func (s *Base) GetLog() Log {
	return s.log
}

func (s *Base) LogError(v ...interface{}) string {
	if s.log == nil {
		return ""
	}

	return s.log.Error(v...)
}

func (s *Base) LogWarning(v ...interface{}) string {
	if s.log == nil {
		return ""
	}

	return s.log.Warning(v...)
}

func (s *Base) LogInfo(v ...interface{}) string {
	if s.log == nil {
		return ""
	}

	return s.log.Info(v...)
}

func (s *Base) LogTrace(v ...interface{}) string {
	if s.log == nil {
		return ""
	}

	return s.log.Trace(v...)
}

func (s *Base) LogDebug(v ...interface{}) string {
	if s.log == nil {
		return ""
	}

	return s.log.Debug(v...)
}
