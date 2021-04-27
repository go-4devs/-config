package config

import "context"

type Key struct {
	Name      string
	AppName   string
	Namespace string
}

type KeyFactory func(ctx context.Context, key Key) string

func (k Key) String() string {
	return k.Namespace + "_" + k.AppName + "_" + k.Name
}
