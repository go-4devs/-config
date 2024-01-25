package definition

import (
	"fmt"
)

func New() Definition {
	return Definition{}
}

type Definition struct {
	options Options
}

func (d *Definition) Add(opts ...Option) *Definition {
	d.options = append(d.options, opts...)

	return d
}

func (d *Definition) View(handle func(Option) error) error {
	for idx, opt := range d.options {
		if err := handle(opt); err != nil {
			return fmt.Errorf("%s[%d]:%w", opt.Kind(), idx, err)
		}
	}

	return nil
}
