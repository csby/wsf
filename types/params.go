package types

type Params interface {
	ByName(name string) string
}
