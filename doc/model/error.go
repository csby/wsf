package model

type Error struct {
	Code    int    `json:"code" note:"错误代码"`
	Summary string `json:"summary" note:"错误描述"`
}

type ErrorSlice []*Error

func (s ErrorSlice) Len() int {
	return len(s)
}

func (s ErrorSlice) Less(i, j int) bool {
	if s[i].Code < s[j].Code {
		return true
	}

	return false
}

func (s ErrorSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
