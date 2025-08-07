package redis

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
)

func TestRealSetNX(t *testing.T) {

	ctx := context.TODO()
	uri := NewUserRedis(
		WithRdsAddr(`127.0.0.1`+":"+`6390`),
		WithRdsPass(`my_redis_password`),
		WithRdsDBNum(2),
	)

	for range 2 {
		err := uri.SetPattern("abc").SetKey("111").SetValue("testValue").Lock().SetNX(ctx).Err()
		if err != nil {
			println(err.Error())
		}
	}

}

func TestSetNX(t *testing.T) {

	ctx := context.TODO()
	db, mock := redismock.NewClientMock()

	uri := &UserRedisImpl{
		Redis:   db,
		process: &Process{},
	}

	mock.ExpectSetNX("abc:111:lock", "1", 1000000000).SetVal(true)
	mock.ExpectSetNX("abc:111", "testValue", 0).SetVal(true)
	mock.ExpectGet("abc:111:lock").SetVal("1")
	mock.ExpectDel("abc:111:lock").SetVal(1)

	err := uri.SetPattern("abc").SetKey("111").SetValue("testValue").Lock().SetNX(ctx).Err()
	if err != nil {
		t.Error(err)
	}
}

// BenchmarkSetNX-12
// 10000           1398112 ns/op            1217 B/op         28 allocs/op
// BenchmarkSetNX-12 (with lock)
// 10000            690767 ns/op            1091 B/op         24 allocs/op
// BenchmarkSetNX-12 (without lock)
// 21266            117751 ns/op             336 B/op          7 allocs/op
func BenchmarkSetNX(b *testing.B) {
	ctx := context.TODO()
	db, mock := redismock.NewClientMock()

	uri := &UserRedisImpl{
		Redis:   db,
		process: &Process{},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		mock.ExpectSetNX("abc:111:lock", strconv.FormatUint(uint64(i)+1, 10), 1000*time.Millisecond).SetVal(true)
		mock.ExpectSetNX("abc:111", "testValue", 0).SetVal(true)
		mock.ExpectGet("abc:111:lock").SetVal(strconv.FormatUint(uint64(i)+1, 10))
		mock.ExpectDel("abc:111:lock").SetVal(1)

		b.StartTimer()

		err := uri.SetPattern("abc").SetKey("111").SetValue("testValue").Lock().SetNX(ctx).Err()
		if err != nil {
			b.Error(err)
		}
	}
}
