package watcher_test

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gitoa.ru/go-4devs/config"
	"gitoa.ru/go-4devs/config/provider/watcher"
	"gitoa.ru/go-4devs/config/value"
)

type provider struct {
	cnt int32
}

func (p *provider) Name() string {
	return "test"
}

func (p *provider) Read(context.Context, config.Key) (config.Variable, error) {
	p.cnt++

	return config.Variable{
		Name:     "tmpname",
		Provider: p.Name(),
		Value:    value.JString(fmt.Sprint(p.cnt)),
	}, nil
}

func TestWatcher(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	prov := &provider{}

	w := watcher.New(time.Second, prov)
	wg := sync.WaitGroup{}
	wg.Add(2)

	var cnt int32

	err := w.Watch(
		ctx,
		config.Key{Name: "tmpname"},
		func(ctx context.Context, oldVar, newVar config.Variable) {
			atomic.AddInt32(&cnt, 1)
			wg.Done()
		},
	)
	require.NoError(t, err)
	wg.Wait()

	require.Equal(t, int32(2), cnt)
}
