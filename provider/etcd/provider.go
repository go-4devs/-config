package etcd

import (
	"context"
	"fmt"
	"log"

	"gitoa.ru/go-4devs/config"
	"gitoa.ru/go-4devs/config/key"
	"gitoa.ru/go-4devs/config/value"
	pb "go.etcd.io/etcd/api/v3/mvccpb"
	client "go.etcd.io/etcd/client/v3"
)

var (
	_ config.Provider      = (*Provider)(nil)
	_ config.WatchProvider = (*Provider)(nil)
)

type Client interface {
	client.KV
	client.Watcher
}

func NewProvider(client Client) *Provider {
	p := Provider{
		client: client,
		key:    key.NsAppName("/"),
	}

	return &p
}

type Provider struct {
	client Client
	key    config.KeyFactory
}

func (p *Provider) IsSupport(ctx context.Context, key config.Key) bool {
	return p.key(ctx, key) != ""
}

func (p *Provider) Name() string {
	return "etcd"
}

func (p *Provider) Read(ctx context.Context, key config.Key) (config.Variable, error) {
	k := p.key(ctx, key)

	resp, err := p.client.Get(ctx, k, client.WithPrefix())
	if err != nil {
		return config.Variable{}, fmt.Errorf("%w: key:%s, prov:%s", err, k, p.Name())
	}

	val, err := p.resolve(k, resp.Kvs)
	if err != nil {
		return config.Variable{}, fmt.Errorf("%w: key:%s, prov:%s", err, k, p.Name())
	}

	return val, nil
}

func (p *Provider) Watch(ctx context.Context, key config.Key, callback config.WatchCallback) error {
	go func(ctx context.Context, key string, callback config.WatchCallback) {
		watch := p.client.Watch(ctx, key, client.WithPrevKV(), client.WithPrefix())
		for w := range watch {
			kvs, olds := p.getEventKvs(w.Events)
			if len(kvs) > 0 {
				newVar, _ := p.resolve(key, kvs)
				oldVar, _ := p.resolve(key, olds)
				callback(ctx, oldVar, newVar)
			}
		}
	}(ctx, p.key(ctx, key), callback)

	return nil
}

func (p *Provider) getEventKvs(events []*client.Event) ([]*pb.KeyValue, []*pb.KeyValue) {
	kvs := make([]*pb.KeyValue, 0, len(events))
	old := make([]*pb.KeyValue, 0, len(events))

	for i := range events {
		kvs = append(kvs, events[i].Kv)
		old = append(old, events[i].PrevKv)
		log.Println(events[i].Type)
	}

	return kvs, old
}

func (p *Provider) resolve(key string, kvs []*pb.KeyValue) (config.Variable, error) {
	for _, kv := range kvs {
		switch {
		case kv == nil:
			return config.Variable{
				Name:     key,
				Provider: p.Name(),
			}, nil
		case string(kv.Key) == key:
			return config.Variable{
				Value:    value.JBytes(kv.Value),
				Name:     key,
				Provider: p.Name(),
			}, nil
		}
	}

	return config.Variable{
		Name:     key,
		Provider: p.Name(),
	}, config.ErrVariableNotFound
}
