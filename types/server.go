package types

type Server interface {
	Interactive() bool
	Run() error
	Restart() error
	Start() error
	Stop() error
	Install() error
	Uninstall() error
	Status() (ServerStatus, error)
}

type ServerStatus byte

const (
	ServerStatusUnknown ServerStatus = iota
	ServerStatusRunning
	ServerStatusStopped
)
