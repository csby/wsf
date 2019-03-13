package types

import "os/exec"

type Module struct {
	Path string
}

func (s *Module) Type() string {
	return s.args("--type")
}

func (s *Module) Name() string {
	return s.args("--module")
}

func (s *Module) Version() string {
	return s.args("--version")
}

func (s *Module) Remark() string {
	return s.args("--remark")
}

func (s *Module) args(name string) string {
	var value = ""
	out, err := exec.Command(s.Path, name).Output()
	if err == nil {
		value = string(out[:])
	}

	return value
}
