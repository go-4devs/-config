package etcd

import (
	"context"
	"fmt"
	"strings"

	"gitoa.ru/go-4devs/config"
	"gitoa.ru/go-4devs/config/value"
	pb "go.etcd.io/etcd/api/v3/mvccpb"
	client "go.etcd.io/etcd/client/v3"
)

const (
	Name      = "etcd"
	Separator = "/"
)

var (
	_ config.Provider      = (*Provider)(nil)
	_ config.WatchProvider = (*Provider)(nil)
)

type Client interface {
	client.KV
	client.Watcher
}

func New(namespace, appName string, client Client) *Provider {
	prov := Provider{
		client: client,
		key: func(s ...string) string {
			return strings.Join(s, Separator)
		},
		name:   Name,
		prefix: namespace + Separator + appName,
	}

	return &prov
}

type Provider struct {
	client Client
	key    func(...string) string
	name   string
	prefix string
	log    func(context.Context, string, ...any)
}

func (p *Provider) Name() string {
	return p.name
}

func (p *Provider) Key(s []string) string {
	return p.prefix + Separator + p.key(s...)
}

func (p *Provider) Value(ctx context.Context, path ...string) (config.Value, error) {
	name := p.Key(path)

	resp, err := p.client.Get(ctx, name, client.WithPrefix())
	if err != nil {
		return nil, fmt.Errorf("%w: key:%s, prov:%s", err, name, p.Name())
	}

	val, err := p.resolve(name, resp.Kvs)
	if err != nil {
		return nil, fmt.Errorf("%w: key:%s, prov:%s", err, name, p.Name())
	}

	return val, nil
}

func (p *Provider) Watch(ctx context.Context, callback config.WatchCallback, path ...string) error {
	go func(ctx context.Context, key string, callback config.WatchCallback) {
		watch := p.client.Watch(ctx, key, client.WithPrevKV(), client.WithPrefix())
		for w := range watch {
			kvs, olds := p.getEventKvs(w.Events)
			if len(kvs) > 0 {
				newVar, _ := p.resolve(key, kvs)
				oldVar, _ := p.resolve(key, olds)

				if err := callback(ctx, oldVar, newVar); err != nil {
					p.log(ctx, "watch callback[%v] %v:%v", p.Name(), path, err)
				}
			}
		}
	}(ctx, p.Key(path), callback)

	return nil
}

func (p *Provider) getEventKvs(events []*client.Event) ([]*pb.KeyValue, []*pb.KeyValue) {
	kvs := make([]*pb.KeyValue, 0, len(events))
	old := make([]*pb.KeyValue, 0, len(events))

	for i := range events {
		kvs = append(kvs, events[i].Kv)
		old = append(old, events[i].PrevKv)
	}

	return kvs, old
}

func (p *Provider) resolve(key string, kvs []*pb.KeyValue) (config.Value, error) {
	for _, kv := range kvs {
		switch {
		case kv == nil:
			return nil, nil
		case string(kv.Key) == key:
			return value.JBytes(kv.Value), nil
		}
	}

	return nil, fmt.Errorf("%w: name %s", config.ErrValueNotFound, key)
}
