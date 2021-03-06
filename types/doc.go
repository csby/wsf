package types

const (
	ContentTypeJson = "application/json"
)

type Function interface {
	SetNote(v string)
	SetTokenType(v int)
	SetInputContentType(v string)
	AddInputHeader(required bool, name, note, defaultValue string, optionValues ...string)
	ClearInputHeader()
	AddInputQuery(required bool, name, note, defaultValue string, optionValues ...string)
	RemoveInputQuery(name string)
	AddInputForm(required bool, key, note string, valueKind int, defaultValue interface{})
	RemoveInputForm(key string)
	SetInputExample(v interface{})
	AddOutputHeader(name, value string)
	ClearOutputHeader()
	SetOutputExample(v interface{})
	SetOutputDataExample(v interface{})
	AddOutputError(err ErrorCode)
	AddOutputErrorCustom(code int, summary string)
}

type Catalog interface {
	AddChild(name string) Catalog
	AddFunction(method string, path HttpPath, name string) Function
}

type Doc interface {
	Enable() bool
	AddCatalog(name string) Catalog
	Catalogs() interface{}
	Function(id, schema, host string) (interface{}, error)
	OnFunctionReady(f func(index int, method, path, name string))
	TokenUI(id string) (interface{}, error)
	TokenCreate(id string, items []TokenAuth, a Assistant) (string, ErrorCode, error)
}
