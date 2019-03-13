package types

type Log interface {
	Error(v ...interface{}) string
	Warning(v ...interface{}) string
	Info(v ...interface{}) string
	Trace(v ...interface{}) string
	Debug(v ...interface{}) string
}
