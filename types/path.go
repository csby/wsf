package types

import "strings"

type Path struct {
	Prefix string
}

func (s *Path) Path(path string) string {
	sb := &strings.Builder{}
	sb.WriteString(s.Prefix)
	sb.WriteString(path)

	return sb.String()
}
