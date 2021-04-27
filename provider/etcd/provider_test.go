package etcd_test

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitoa.ru/go-4devs/config"
	"gitoa.ru/go-4devs/config/provider/etcd"
	"gitoa.ru/go-4devs/config/test"
)

func TestProvider(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	et, err := test.NewEtcd(ctx)
	require.NoError(t, err)

	provider := etcd.NewProvider(et)
	read := []test.Read{
		test.NewRead("db_dsn", test.DSN),
		test.NewRead("duration", 12*time.Minute),
		test.NewRead("port", 8080),
		test.NewRead("maintain", true),
		test.NewRead("start_at", test.Time("2020-01-02T15:04:05Z")),
		test.NewRead("percent", .064),
		test.NewRead("count", uint(2020)),
		test.NewRead("int64", int64(2021)),
		test.NewRead("uint64", int64(2022)),
		test.NewReadConfig("config"),
	}
	test.Run(t, provider, read)
}

func value(cnt int32) string {
	return fmt.Sprintf("test data: %d", cnt)
}

func TestWatcher(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	key := config.Key{
		AppName:   "config",
		Namespace: "fdevs",
		Name:      "test_watch",
	}

	et, err := test.NewEtcd(ctx)
	require.NoError(t, err)

	defer func() {
		_, err = et.KV.Delete(context.Background(), "fdevs/config/test_watch")
		require.NoError(t, err)
	}()

	var cnt, cnt2 int32

	prov := etcd.NewProvider(et)
	wg := sync.WaitGroup{}
	wg.Add(6)

	watch := func(cnt *int32) func(ctx context.Context, oldVar, newVar config.Variable) {
		return func(ctx context.Context, oldVar, newVar config.Variable) {
			switch *cnt {
			case 0:
				assert.Equal(t, value(*cnt), newVar.Value.String())
				assert.Nil(t, oldVar.Value)
			case 1:
				assert.Equal(t, value(*cnt), newVar.Value.String())
				assert.Equal(t, value(*cnt-1), oldVar.Value.String())
			case 2:
				_, perr := newVar.Value.ParseString()
				assert.NoError(t, perr)
				assert.Equal(t, "", newVar.Value.String())
				assert.Equal(t, value(*cnt-1), oldVar.Value.String())
			default:
				assert.Fail(t, "unexpected watch")
			}

			wg.Done()
			atomic.AddInt32(cnt, 1)
		}
	}

	err = prov.Watch(ctx, key, watch(&cnt))
	err = prov.Watch(ctx, key, watch(&cnt2))
	require.NoError(t, err)

	time.AfterFunc(time.Second, func() {
		_, err = et.KV.Put(ctx, "fdevs/config/test_watch", value(0))
		require.NoError(t, err)
		_, err = et.KV.Put(ctx, "fdevs/config/test_watch", value(1))
		require.NoError(t, err)
		_, err = et.KV.Delete(ctx, "fdevs/config/test_watch")
		require.NoError(t, err)
	})

	time.AfterFunc(time.Second*10, func() {
		assert.Fail(t, "failed watch after 5 sec")
		cancel()
	})

	go func() {
		wg.Wait()
		cancel()
	}()

	<-ctx.Done()
}
