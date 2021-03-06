package model

import (
	"encoding/json"
	"testing"
)

func TestArgument_FromExample(t *testing.T) {
	argument := &Argument{}

	argumentBase := ArgumentBase{
		Data: uint64(50),
	}
	example := argument.FromExample(argumentBase)
	t.Logf("%v", toJson(example))

	argumentStruct := ArgumentStruct{}
	example = argument.FromExample(argumentStruct)
	t.Logf("%v", toJson(example))

	argumentPoint := ArgumentPoint{
		Object: &argumentBase,
	}
	example = argument.FromExample(argumentPoint)
	t.Logf("%v", toJson(example))
}

func TestArgument_FromExample_Anonymous(t *testing.T) {
	argument := &Argument{}

	argsChild := &ArgsChild{
		ArgsParent: ArgsParent{
			ParentId: "12",
		},
		ChildId: uint(50),
	}
	example := argument.FromExample(argsChild)
	t.Logf("%v", toJson(example))

	models := example.ToModel()
	t.Logf("%v", toJson(models))
}

type ArgumentBase struct {
	Name   string      `json:"name" note:"名称(Base)"`
	Value  int         `json:"value" note:"值(base)-int"`
	Point  *float64    `json:"point" note:"值(base)-*float64"`
	Data   interface{} `json:"data" note:"数据"`
	Arrays []string    `json:"arrays" note:"数组"`
}

type ArgumentStruct struct {
	Kind   string       `json:"kind" note:"struct" required:"true"`
	Object ArgumentBase `json:"object"`
}

type ArgumentPoint struct {
	Kind   string        `json:"kind" note:"point" required:"true"`
	Object *ArgumentBase `json:"object"`
}

func toJson(v interface{}) string {
	bytes, err := json.MarshalIndent(v, "", "    ")
	if err != nil {
		return ""
	}

	return string(bytes[:])
}

type ArgsParent struct {
	ParentId string `json:"parentId" required:"true" note:"父ID"`
}

type ArgsChild struct {
	ArgsParent

	ChildId interface{} `json:"childId" required:"false" note:"子ID"`
}
