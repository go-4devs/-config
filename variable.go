package config

type Variable struct {
	Name     string
	Provider string
	Value    Value
}

func (v Variable) IsEquals(n Variable) bool {
	return n.Name == v.Name && n.Provider == v.Provider && n.Value.String() == v.Value.String()
}
