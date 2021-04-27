package yaml

import (
	"context"
	"fmt"
	"io/ioutil"

	"gitoa.ru/go-4devs/config"
	"gopkg.in/yaml.v3"
)

func WithFileKeyFactory(f func(context.Context, config.Key) []string) FileOption {
	return func(p *File) {
		p.key = f
	}
}

type FileOption func(*File)

func NewFile(name string, opts ...FileOption) *File {
	f := File{
		file: name,
		key:  keyFactory,
	}

	for _, opt := range opts {
		opt(&f)
	}

	return &f
}

type File struct {
	file string
	key  func(context.Context, config.Key) []string
}

func (p *File) Name() string {
	return "yaml_file"
}

func (p *File) Read(ctx context.Context, key config.Key) (config.Variable, error) {
	in, err := ioutil.ReadFile(p.file)
	if err != nil {
		return config.Variable{}, fmt.Errorf("yaml_file: read error: %w", err)
	}

	var n yaml.Node
	if err = yaml.Unmarshal(in, &n); err != nil {
		return config.Variable{}, fmt.Errorf("yaml_file: unmarshal error: %w", err)
	}

	data := node{Node: &n}
	k := p.key(ctx, key)

	return data.read(p.Name(), k)
}
