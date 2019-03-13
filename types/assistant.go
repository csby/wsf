package types

import "time"

type Assistant interface {
	// input & output
	GetBody() ([]byte, error)
	GetJson(v interface{}) error
	GetXml(v interface{}) error
	OutputJson(v interface{})
	OutputXml(v interface{})
	Success(data interface{})
	Error(errCode ErrorCode, errDetails ...interface{})

	// service
	CanUpdate() bool
	CanRestart() bool
	Restart() error

	// others
	IsError() bool
	Set(key string, val interface{})
	Get(key string) (interface{}, bool)
	Del(key string) bool
	SetLog(v bool)
	GetLog() bool
	SetInput(v []byte)
	GetInput() []byte
	GetOutput() []byte
	GetParam() []byte
	Method() string
	Schema() string
	Path() string
	RID() uint64
	RIP() string
	EnterTime() time.Time
	LeaveTime() time.Time
	Token() string
	NewGuid() string
}
