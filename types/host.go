package types

type Host interface {
	Run() error
	Close() error
}
