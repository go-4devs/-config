package definition

type Option interface {
	Kind() string
}

type Options []Option

func (s Options) Len() int      { return len(s) }
func (s Options) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

type Params []Param

func (p Params) Get(name string) (any, bool) {
	for _, param := range p {
		if param.Name == name {
			return param.Value, true
		}
	}

	return nil, false
}

type Param struct {
	Name  string
	Value any
}
