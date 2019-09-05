package model

import (
	"fmt"
	"reflect"
	"strings"
)

const (
	tagJson     = "json"
	tagNote     = "note"
	tagRequired = "required"
)

type Argument struct {
	parent *Argument

	Name     string `json:"name"`     // 名称
	Type     string `json:"type"`     // 类型
	Note     string `json:"note"`     // 说明
	Required bool   `json:"required"` // 必填

	Children []*Argument `json:"children"`
}

func (s *Argument) ParentType() string {
	if s.parent == nil {
		return ""
	}

	return s.parent.Type
}

func (s *Argument) FromExample(example interface{}) *Argument {
	if example == nil {
		return nil
	}

	model := &Argument{Children: make([]*Argument, 0)}

	exampleType := reflect.TypeOf(example)
	exampleTypeKind := exampleType.Kind()
	switch exampleTypeKind {
	case reflect.Ptr:
		{
			s.parseExample(reflect.ValueOf(example).Elem(), model)
			break
		}
	case reflect.Interface,
		reflect.Struct,
		reflect.Array,
		reflect.Slice,
		reflect.Bool,
		reflect.String,
		reflect.Float32, reflect.Float64,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		{
			s.parseExample(reflect.ValueOf(example), model)
			break
		}
	default:
		return nil
	}

	return model
}

func (s *Argument) parseExample(v reflect.Value, argument *Argument) {
	if argument == nil {
		return
	}
	if v.Kind() == reflect.Invalid {
		return
	}

	t := v.Type()
	k := t.Kind()
	switch k {
	case reflect.Ptr:
		{
			if argument.Type == "" {
				argument.Type = t.String()
			}
			s.parseExample(v.Elem(), argument)
			break
		}
	case reflect.Interface:
		{
			if argument.Type == "" {
				argument.Type = k.String()
			}
			if v.CanInterface() {
				value := reflect.ValueOf(v.Interface())
				if value.Kind() != reflect.Invalid {
					s.parseExample(value, argument)
				}
			}
			break
		}
	case reflect.Struct:
		{
			if argument.Type == "" {
				argument.Type = t.Name()
			}

			n := v.NumField()
			for i := 0; i < n; i++ {
				valueField := v.Field(i)
				if !valueField.CanInterface() {
					continue
				}

				typeField := t.Field(i)
				if typeField.Anonymous {
					if valueField.CanAddr() {
						s.parseExample(valueField.Addr().Elem(), argument)
					}
				} else {
					child := &Argument{Children: make([]*Argument, 0)}
					child.Name = typeField.Tag.Get(tagJson)
					if child.Name == "" {
						child.Name = typeField.Name
					}
					cns := strings.Split(child.Name, ",")
					if len(cns) > 1 {
						child.Name = cns[0]
					}
					child.Type = valueField.Kind().String()
					//child.Type = valueField.Type().String()
					if typeField.Tag.Get(tagRequired) == "true" {
						child.Required = true
					}
					child.Note = typeField.Tag.Get(tagNote)
					child.parent = argument
					argument.Children = append(argument.Children, child)

					value := reflect.ValueOf(valueField.Interface())
					if value.Kind() != reflect.Invalid {
						child.Type = value.Type().Name()
						s.parseExample(value, child)
						if child.Type == "" {
							child.Type = valueField.Type().String()
						}
					}
				}
			}
			break
		}
	case reflect.Array:
		{
			break
		}
	case reflect.Slice:
		{
			st := t.Elem()
			stk := st.Kind()

			var ste reflect.Type = nil
			if stk == reflect.Ptr {
				ste = st.Elem()
			} else {
				ste = st
			}
			if ste != nil {
				argument.Type = fmt.Sprintf("%s[]", ste.Name())

				if ste.Kind() == reflect.Struct && argument.ParentType() != ste.Name() {
					stet := reflect.New(ste)
					child := &Argument{Children: make([]*Argument, 0)}
					child.Type = ste.Name()
					argument.Children = append(argument.Children, child)
					s.parseExample(stet.Elem(), child)
				}
			} else {
				argument.Type = fmt.Sprintf("%s[]", stk.String())
			}

			if argument.Type == "[]" {
				argument.Type = fmt.Sprintf("%s[]", stk.String())
			}

			break
		}
	case reflect.Bool,
		reflect.String,
		reflect.Float32, reflect.Float64,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		{
			argument.Type = t.Name()
			break
		}
	default:
		{
			return
		}
	}

}

func (s *Argument) ToModel() []*Model {
	results := make([]*Model, 0)
	model := &Model{
		Name:     s.typeToKind(s.Type),
		Children: make([]*Item, 0),
	}
	types := make(map[string]string)
	s.toModel(&results, model, s.Children, types)

	return results
}

func (s *Argument) toModel(results *[]*Model, model *Model, children []*Argument, types map[string]string) {
	childCount := len(children)
	if childCount < 1 {
		return
	}

	_, ok := types[model.Name]
	if ok {
		return
	}
	types[model.Name] = ""
	if !strings.Contains(model.Name, "[]") {
		*results = append(*results, model)
	}

	for childIndex := 0; childIndex < childCount; childIndex++ {
		child := children[childIndex]
		if child == nil {
			continue
		}

		item := &Item{}
		child.copyTo(item)
		model.Children = append(model.Children, item)
		if len(child.Children) > 0 {
			childModel := &Model{
				Name:     s.typeToKind(child.Type),
				Children: make([]*Item, 0),
			}
			s.toModel(results, childModel, child.Children, types)
		}
	}
}

func (s *Argument) typeToKind(name string) string {
	return strings.ReplaceAll(strings.ReplaceAll(name, "*", ""), "{}", "")
}

func (s *Argument) copyTo(target *Item) {
	if target == nil {
		return
	}

	target.Name = s.Name
	target.Type = s.Type
	target.Note = s.Note
	target.Required = s.Required
}
